-- Copyright 2024 The Refactored Authors. All Rights Reserved.
-- Database Performance Optimization Script
-- Run this script to optimize database queries for the refactored API modules

-- ============================================
-- User Table Indexes
-- ============================================

-- Primary lookup indexes
CREATE INDEX IF NOT EXISTS idx_user_owner_name ON "user" (owner, name);
CREATE INDEX IF NOT EXISTS idx_user_owner_email ON "user" (owner, email);
CREATE INDEX IF NOT EXISTS idx_user_owner_phone ON "user" (owner, phone);
CREATE INDEX IF NOT EXISTS idx_user_owner_user_id ON "user" (owner, user_id);

-- Search indexes
CREATE INDEX IF NOT EXISTS idx_user_name ON "user" (name);
CREATE INDEX IF NOT EXISTS idx_user_email ON "user" (email);
CREATE INDEX IF NOT EXISTS idx_user_phone ON "user" (phone);
CREATE INDEX IF NOT EXISTS idx_user_display_name ON "user" (display_name);

-- Filter indexes
CREATE INDEX IF NOT EXISTS idx_user_is_forbidden ON "user" (is_forbidden);
CREATE INDEX IF NOT EXISTS idx_user_is_deleted ON "user" (is_deleted);
CREATE INDEX IF NOT EXISTS idx_user_is_admin ON "user" (is_admin);
CREATE INDEX IF NOT EXISTS idx_user_signup_application ON "user" (signup_application);

-- Sorting indexes
CREATE INDEX IF NOT EXISTS idx_user_created_time ON "user" (created_time);
CREATE INDEX IF NOT EXISTS idx_user_score ON "user" (score);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_user_owner_created_time ON "user" (owner, created_time);
CREATE INDEX IF NOT EXISTS idx_user_owner_score ON "user" (owner, score DESC);

-- ============================================
-- Organization Table Indexes
-- ============================================

-- Primary lookup indexes
CREATE INDEX IF NOT EXISTS idx_organization_owner_name ON organization (owner, name);
CREATE INDEX IF NOT EXISTS idx_organization_name ON organization (name);

-- Hierarchy index
CREATE INDEX IF NOT EXISTS idx_organization_parent_id ON organization (parent_id);

-- Search indexes
CREATE INDEX IF NOT EXISTS idx_organization_display_name ON organization (display_name);
CREATE INDEX IF NOT EXISTS idx_organization_website_url ON organization (website_url);

-- Filter indexes
CREATE INDEX IF NOT EXISTS idx_organization_enable_soft_deletion ON organization (enable_soft_deletion);
CREATE INDEX IF NOT EXISTS idx_organization_is_profile_public ON organization (is_profile_public);

-- Sorting indexes
CREATE INDEX IF NOT EXISTS idx_organization_created_time ON organization (created_time);

-- Composite indexes
CREATE INDEX IF NOT EXISTS idx_organization_owner_created_time ON organization (owner, created_time);

-- ============================================
-- Application Table Indexes
-- ============================================

-- Primary lookup indexes
CREATE INDEX IF NOT EXISTS idx_application_owner_name ON application (owner, name);
CREATE INDEX IF NOT EXISTS idx_application_client_id ON application (client_id);

-- Search indexes
CREATE INDEX IF NOT EXISTS idx_application_name ON application (name);
CREATE INDEX IF NOT EXISTS idx_application_display_name ON application (display_name);

-- Filter indexes
CREATE INDEX IF NOT EXISTS idx_application_organization ON application (organization);
CREATE INDEX IF NOT EXISTS idx_application_is_shared ON application (is_shared);
CREATE INDEX IF NOT EXISTS idx_application_enable_sign_up ON application (enable_sign_up);

-- Sorting indexes
CREATE INDEX IF NOT EXISTS idx_application_created_time ON application (created_time);

-- Composite indexes
CREATE INDEX IF NOT EXISTS idx_application_owner_created_time ON application (owner, created_time);

-- ============================================
-- Group Table Indexes
-- ============================================

CREATE INDEX IF NOT EXISTS idx_group_owner_name ON "group" (owner, name);
CREATE INDEX IF NOT EXISTS idx_group_parent_id ON "group" (parent_id);
CREATE INDEX IF NOT EXISTS idx_group_created_time ON "group" (created_time);

-- ============================================
-- Permission Table Indexes
-- ============================================

CREATE INDEX IF NOT EXISTS idx_permission_owner_name ON permission (owner, name);
CREATE INDEX IF NOT EXISTS idx_permission_resource_type ON permission (resource_type);
CREATE INDEX IF NOT EXISTS idx_permission_created_time ON permission (created_time);

-- ============================================
-- Role Table Indexes
-- ============================================

CREATE INDEX IF NOT EXISTS idx_role_owner_name ON role (owner, name);
CREATE INDEX IF NOT EXISTS idx_role_created_time ON role (created_time);

-- ============================================
-- Record Table Indexes
-- ============================================

CREATE INDEX IF NOT EXISTS idx_record_owner_name ON record (owner, name);
CREATE INDEX IF NOT EXISTS idx_record_created_time ON record (created_time);
CREATE INDEX IF NOT EXISTS idx_record_organization ON record (organization);

-- ============================================
-- Full Text Search (PostgreSQL only)
-- ============================================

-- For PostgreSQL, enable full text search
-- CREATE INDEX IF NOT EXISTS idx_user_fulltext ON "user" USING gin(to_tsvector('english', name || ' ' || COALESCE(display_name, '') || ' ' || COALESCE(email, '')));
-- CREATE INDEX IF NOT EXISTS idx_organization_fulltext ON organization USING gin(to_tsvector('english', name || ' ' || COALESCE(display_name, '') || ' ' || COALESCE(website_url, '')));
-- CREATE INDEX IF NOT EXISTS idx_application_fulltext ON application USING gin(to_tsvector('english', name || ' ' || COALESCE(display_name, '') || ' ' || COALESCE(description, '')));

-- ============================================
-- Table Statistics Update
-- ============================================

-- Update table statistics for query optimizer
ANALYZE "user";
ANALYZE organization;
ANALYZE application;
ANALYZE "group";
ANALYZE permission;
ANALYZE role;
ANALYZE record;

-- ============================================
-- Query Performance Monitoring
-- ============================================

-- Enable query logging (PostgreSQL)
-- ALTER SYSTEM SET log_min_duration_statement = 1000;  -- Log queries taking more than 1 second
-- SELECT pg_reload_conf();

-- Enable slow query log (MySQL)
-- SET GLOBAL slow_query_log = 'ON';
-- SET GLOBAL long_query_time = 1;
