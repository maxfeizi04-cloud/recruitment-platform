package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool
var ctx = context.Background()

func main() {
	dsn := getEnv("DATABASE_URL", "postgres://devuser:devpass@192.168.198.133:5432/recruitment?sslmode=disable")

	var err error
	pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	fmt.Println("🌱 开始生成种子数据...\n")

	// 清理旧数据
	pool.Exec(ctx, `DELETE FROM interview_invitations`)
	pool.Exec(ctx, `DELETE FROM job_applications`)
	pool.Exec(ctx, `DELETE FROM jobs`)
	pool.Exec(ctx, `DELETE FROM resumes`)
	pool.Exec(ctx, `DELETE FROM hr_certifications`)
	pool.Exec(ctx, `DELETE FROM users`)
	fmt.Println("✅ 清理旧数据完成\n")

	// 创建用户
	hrIDs := createHRUsers()
	candidateIDs := createCandidates()
	fmt.Printf("✅ 创建 %d 个HR + %d 个求职者\n\n", len(hrIDs), len(candidateIDs))

	// 创建简历
	resumeIDs := createResumes(candidateIDs)
	fmt.Printf("✅ 创建 %d 份简历\n\n", len(resumeIDs))

	// 创建职位
	jobIDs := createJobs(hrIDs)
	fmt.Printf("✅ 创建 %d 个职位\n\n", len(jobIDs))

	// 创建投递记录
	appIDs := createApplications(candidateIDs, resumeIDs, jobIDs)
	fmt.Printf("✅ 创建 %d 条投递记录\n\n", len(appIDs))

	// 创建面试邀约
	invIDs := createInterviews(hrIDs, candidateIDs, appIDs, jobIDs)
	fmt.Printf("✅ 创建 %d 条面试邀约\n\n", len(invIDs))

	fmt.Println("🎉 种子数据生成完毕！")
	fmt.Println()
	fmt.Println("测试账号（验证码: 123456）：")
	fmt.Println("  求职者: 13800000001  验证码存入 Redis: SET sms:code:13800000001 123456")
	fmt.Println("  HR:    13900000001  验证码存入 Redis: SET sms:code:13900000001 123456")
	fmt.Println("  管理员: 13700000001  验证码存入 Redis: SET sms:code:13700000001 123456")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ── Users ──

func createHRUsers() []uuid.UUID {
	hrs := []struct{ phone, name, company, position string }{
		{"13900000001", "李经理", "字节跳动", "HR总监"},
		{"13900000002", "王主管", "腾讯科技", "招聘经理"},
		{"13900000003", "张HR", "阿里巴巴", "HRBP"},
	}
	var ids []uuid.UUID
	for _, h := range hrs {
		var id uuid.UUID
		pool.QueryRow(ctx, `INSERT INTO users (phone, role, name) VALUES ($1, 'hr', $2) ON CONFLICT (phone) DO UPDATE SET name=$2 RETURNING id`, h.phone, h.name).Scan(&id)
		pool.Exec(ctx, `INSERT INTO hr_certifications (user_id, company_name, position, status) VALUES ($1, $2, $3, 'approved') ON CONFLICT (user_id) DO UPDATE SET company_name=$2, position=$3`, id, h.company, h.position)
		ids = append(ids, id)
		fmt.Printf("  HR: %s (%s - %s)\n", h.name, h.company, h.position)
	}
	return ids
}

func createCandidates() []uuid.UUID {
	candidates := []struct{ phone, name string }{
		{"13800000001", "张一一"},
		{"13800000002", "陈小明"},
		{"13800000003", "刘技术"},
		{"13800000004", "赵产品"},
		{"13800000005", "孙设计"},
		{"13800000006", "周前端"},
		{"13800000007", "吴后端"},
		{"13800000008", "郑数据"},
		{"13800000009", "王运维"},
		{"13800000010", "钱测试"},
	}
	var ids []uuid.UUID
	for _, c := range candidates {
		var id uuid.UUID
		pool.QueryRow(ctx, `INSERT INTO users (phone, role, name) VALUES ($1, 'candidate', $2) ON CONFLICT (phone) DO UPDATE SET name=$2 RETURNING id`, c.phone, c.name).Scan(&id)
		ids = append(ids, id)
	}
	// Admin
	pool.Exec(ctx, `INSERT INTO users (phone, role, name) VALUES ('13700000001', 'admin', '管理员') ON CONFLICT (phone) DO NOTHING`)

	return ids
}

// ── Resumes ──

func createResumes(candidateIDs []uuid.UUID) []uuid.UUID {
	type resumeData struct {
		userIdx int
		title   string
		skills  string
	}
	data := []resumeData{
		{0, "全栈工程师简历", `{"skills":["React","TypeScript","Go","PostgreSQL","Docker"]}`},
		{1, "高级前端工程师简历", `{"skills":["React","Vue","TypeScript","Webpack","CSS"]}`},
		{2, "后端开发工程师简历", `{"skills":["Go","Java","Python","K8s","Microservices"]}`},
		{3, "产品经理简历", `{"skills":["产品设计","用户研究","数据分析","Axure","Figma"]}`},
		{4, "UI设计专家简历", `{"skills":["Figma","Sketch","用户体验","交互设计","品牌设计"]}`},
	}
	var ids []uuid.UUID
	for _, d := range data {
		var id uuid.UUID
		userID := candidateIDs[d.userIdx]
		pool.QueryRow(ctx,
			`INSERT INTO resumes (user_id, title, content, is_default) VALUES ($1, $2, $3::jsonb, true) RETURNING id`,
			userID, d.title, d.skills,
		).Scan(&id)
		ids = append(ids, id)
		fmt.Printf("  简历: %s (用户 %d)\n", d.title, d.userIdx+1)
	}
	return ids
}

// ── Jobs ──

func createJobs(hrIDs []uuid.UUID) []uuid.UUID {
	type jobData struct {
		hrIdx       int
		title       string
		description string
		skills      string
		salary      string
		location    string
	}
	jobs := []jobData{
		{0, "高级前端工程师", "负责核心产品前端架构设计与开发，优化页面性能", `{"skills":["React","TypeScript","Webpack"]}`, `{"min":25000,"max":40000,"period":"monthly"}`, `{"province":"北京","city":"北京","district":"朝阳区"}`},
		{0, "资深后端开发", "负责微服务架构设计及核心模块开发", `{"skills":["Go","Docker","K8s","PostgreSQL"]}`, `{"min":30000,"max":50000,"period":"monthly"}`, `{"province":"北京","city":"北京","district":"海淀区"}`},
		{0, "产品经理", "负责B端产品规划与需求管理", `{"skills":["产品设计","数据分析","Axure"]}`, `{"min":20000,"max":35000,"period":"monthly"}`, `{"province":"北京","city":"北京","district":"朝阳区"}`},
		{0, "DevOps工程师", "负责CI/CD流水线建设与云原生基础设施", `{"skills":["K8s","Docker","Terraform","Jenkins"]}`, `{"min":28000,"max":45000,"period":"monthly"}`, `{"province":"北京","city":"北京","district":"海淀区"}`},
		{1, "全栈开发工程师", "参与微信生态产品全栈开发", `{"skills":["React","Node.js","Go","MySQL"]}`, `{"min":22000,"max":38000,"period":"monthly"}`, `{"province":"深圳","city":"深圳","district":"南山区"}`},
		{1, "iOS开发工程师", "负责短视频App iOS客户端开发", `{"skills":["Swift","SwiftUI","Combine","CoreData"]}`, `{"min":25000,"max":42000,"period":"monthly"}`, `{"province":"深圳","city":"深圳","district":"南山区"}`},
		{1, "游戏后台开发", "负责大型多人在线游戏后台开发", `{"skills":["C++","Go","Redis","Linux"]}`, `{"min":30000,"max":55000,"period":"monthly"}`, `{"province":"深圳","city":"深圳","district":"南山区"}`},
		{1, "测试开发工程师", "负责质量平台建设和自动化测试框架", `{"skills":["Python","Go","Selenium","JMeter"]}`, `{"min":20000,"max":35000,"period":"monthly"}`, `{"province":"深圳","city":"深圳","district":"宝安区"}`},
		{2, "Java开发工程师", "负责电商交易核心链路开发", `{"skills":["Java","Spring","MyBatis","MySQL","Redis"]}`, `{"min":25000,"max":40000,"period":"monthly"}`, `{"province":"杭州","city":"杭州","district":"余杭区"}`},
		{2, "算法工程师", "负责搜索推荐算法优化", `{"skills":["Python","TensorFlow","PyTorch","Spark"]}`, `{"min":35000,"max":60000,"period":"monthly"}`, `{"province":"杭州","city":"杭州","district":"余杭区"}`},
		{2, "数据工程师", "负责数据仓库建设和ETL流程优化", `{"skills":["SQL","Python","Spark","Flink","Hive"]}`, `{"min":28000,"max":45000,"period":"monthly"}`, `{"province":"杭州","city":"杭州","district":"西湖区"}`},
		{2, "安全工程师", "负责应用安全和渗透测试", `{"skills":["Web安全","渗透测试","Python","OWASP"]}`, `{"min":30000,"max":50000,"period":"monthly"}`, `{"province":"杭州","city":"杭州","district":"余杭区"}`},
		{0, "UI/UX设计师", "负责产品界面视觉设计与用户体验优化", `{"skills":["Figma","Sketch","用户体验","交互设计"]}`, `{"min":18000,"max":30000,"period":"monthly"}`, `{"province":"上海","city":"上海","district":"浦东新区"}`},
		{1, "数据分析师", "负责业务数据分析和BI报表建设", `{"skills":["SQL","Python","Tableau","Excel"]}`, `{"min":18000,"max":32000,"period":"monthly"}`, `{"province":"上海","city":"上海","district":"静安区"}`},
		{2, "技术经理", "负责研发团队管理和技术规划", `{"skills":["团队管理","技术规划","Agile","系统架构"]}`, `{"min":40000,"max":65000,"period":"monthly"}`, `{"province":"上海","city":"上海","district":"徐汇区"}`},
		{0, "Android工程师", "负责短视频App安卓端开发", `{"skills":["Kotlin","Jetpack","Compose","Coroutines"]}`, `{"min":22000,"max":38000,"period":"monthly"}`, `{"province":"广州","city":"广州","district":"天河区"}`},
		{1, "运维工程师", "负责线上服务稳定性保障", `{"skills":["Linux","K8s","Prometheus","Grafana","Ansible"]}`, `{"min":20000,"max":35000,"period":"monthly"}`, `{"province":"广州","city":"广州","district":"海珠区"}`},
		{2, "技术文档工程师", "负责技术文档撰写和API文档维护", `{"skills":["Markdown","OpenAPI","技术写作"]}`, `{"min":15000,"max":25000,"period":"monthly"}`, `{"province":"广州","city":"广州","district":"越秀区"}`},
		{0, "区块链开发", "负责Web3 dApp智能合约开发", `{"skills":["Solidity","Ethereum","Web3.js","Go"]}`, `{"min":35000,"max":60000,"period":"monthly"}`, `{"province":"杭州","city":"杭州","district":"滨江区"}`},
		{1, "AIGC工程师", "负责大模型应用开发与Prompt工程", `{"skills":["Python","LangChain","OpenAI","VectorDB"]}`, `{"min":40000,"max":70000,"period":"monthly"}`, `{"province":"北京","city":"北京","district":"海淀区"}`},
	}
	var ids []uuid.UUID
	for _, j := range jobs {
		var id uuid.UUID
		pool.QueryRow(ctx,
			`INSERT INTO jobs (hr_user_id, title, description, requirements, salary_range, location, status)
			 VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6::jsonb, 'active') RETURNING id`,
			hrIDs[j.hrIdx], j.title, j.description, j.skills, j.salary, j.location,
		).Scan(&id)
		ids = append(ids, id)
	}
	fmt.Printf("  职位: %d 个（分布在北上广深杭）\n", len(jobs))
	return ids
}

// ── Applications ──

func createApplications(candidateIDs, resumeIDs, jobIDs []uuid.UUID) []uuid.UUID {
	apps := []struct{ candIdx, resumeIdx, jobIdx int }{
		{0, 0, 0}, // 张一一 → 高级前端
		{0, 0, 9}, // 张一一 → Java开发
		{1, 1, 0}, // 陈小明 → 高级前端
		{2, 2, 1}, // 刘技术 → 资深后端
		{3, 3, 2}, // 赵产品 → 产品经理
	}
	var ids []uuid.UUID
	for _, a := range apps {
		var id uuid.UUID
		err := pool.QueryRow(ctx,
			`INSERT INTO job_applications (job_id, user_id, resume_id, status)
			 VALUES ($1, $2, $3, 'pending')
			 ON CONFLICT (job_id, user_id) DO NOTHING
			 RETURNING id`,
			jobIDs[a.jobIdx], candidateIDs[a.candIdx], resumeIDs[a.resumeIdx],
		).Scan(&id)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// ── Interviews ──

func createInterviews(hrIDs, candidateIDs, appIDs, jobIDs []uuid.UUID) []uuid.UUID {
	invitations := []struct{ hrIdx, candIdx, jobIdx int }{
		{0, 0, 0}, // 字节HR → 张一一 → 高级前端
		{1, 1, 0}, // 腾讯HR → 陈小明 → 高级前端
	}
	var ids []uuid.UUID
	addresses := []string{
		`{"province":"北京","city":"北京","district":"朝阳区","detail":"望京SOHO T1 20层","lat":39.997,"lng":116.479,"formatted":"北京市朝阳区望京SOHO T1"}`,
		`{"province":"深圳","city":"深圳","district":"南山区","detail":"科技园南路88号","lat":22.543,"lng":113.954,"formatted":"深圳市南山区科技园南路88号"}`,
	}
	for i, inv := range invitations {
		var id uuid.UUID
		// Find the application ID for this candidate+job combination
		var appID uuid.UUID
		err := pool.QueryRow(ctx,
			`SELECT id FROM job_applications WHERE user_id=$1 AND job_id=$2 LIMIT 1`,
			candidateIDs[inv.candIdx], jobIDs[inv.jobIdx],
		).Scan(&appID)
		if err != nil {
			continue
		}
		pool.QueryRow(ctx,
			`INSERT INTO interview_invitations (job_application_id, hr_user_id, candidate_user_id, scheduled_at, company_address, contact_name, contact_phone, notes, status)
			 VALUES ($1, $2, $3, '2026-07-15 10:00:00+08', $4::jsonb, 'HR经理', '13800000000', '请携带简历和作品集', 'pending') RETURNING id`,
			appID, hrIDs[inv.hrIdx], candidateIDs[inv.candIdx], addresses[i],
		).Scan(&id)
		ids = append(ids, id)
	}
	return ids
}
