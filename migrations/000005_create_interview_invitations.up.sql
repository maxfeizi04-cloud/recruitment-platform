CREATE TABLE interview_invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_application_id UUID NOT NULL REFERENCES job_applications(id) ON DELETE CASCADE,
    hr_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    candidate_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    company_address JSONB NOT NULL DEFAULT '{}',
    contact_name VARCHAR(50) NOT NULL DEFAULT '',
    contact_phone VARCHAR(20) NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    attachment_urls TEXT[] NOT NULL DEFAULT '{}',
    status VARCHAR(15) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined', 'reschedule', 'confirmed')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_interview_invitations_hr ON interview_invitations (hr_user_id);
CREATE INDEX idx_interview_invitations_candidate ON interview_invitations (candidate_user_id);
CREATE INDEX idx_interview_invitations_status ON interview_invitations (status);
