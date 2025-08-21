-- Database initialization script for Product Management API

-- Create database if not exists (this line is commented out because PostgreSQL 
-- automatically creates the database specified in POSTGRES_DB)
-- CREATE DATABASE IF NOT EXISTS product_management;

-- Connect to the database
\c product_management;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types if needed
-- CREATE TYPE user_role AS ENUM ('user', 'admin');

-- Tables will be created automatically by GORM migrations
-- This script is mainly for any initial data or database setup

-- Insert default admin user (password: admin123)
-- This will be handled by the application or can be uncommented if needed
/*
INSERT INTO users (
    email, 
    username, 
    password, 
    first_name, 
    last_name, 
    is_active, 
    is_admin, 
    created_at, 
    updated_at
) VALUES (
    'admin@example.com',
    'admin',
    '$2a$10$8kVVmvRME0mILIR4CuDJruMvEt/XKSNfr7w8vULI9R5BN9VK8bN7e', -- admin123
    'System',
    'Administrator',
    true,
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;
*/

-- Insert sample product categories
-- These will be created dynamically by the API

-- Create indexes for better performance (GORM will handle most of these)
-- Additional custom indexes can be added here

-- Sample data for development (optional)
-- Uncomment the following lines to insert sample data

/*
-- Sample products
INSERT INTO products (
    name, 
    description, 
    price, 
    stock, 
    category, 
    image_url, 
    is_active, 
    created_at, 
    updated_at
) VALUES 
    ('Sample Product 1', 'This is a sample product for testing', 29.99, 100, 'Electronics', 'https://example.com/image1.jpg', true, NOW(), NOW()),
    ('Sample Product 2', 'Another sample product', 49.99, 50, 'Clothing', 'https://example.com/image2.jpg', true, NOW(), NOW()),
    ('Sample Product 3', 'Third sample product', 19.99, 200, 'Books', 'https://example.com/image3.jpg', true, NOW(), NOW())
ON CONFLICT DO NOTHING;
*/

-- Grant necessary permissions (if using specific database users)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_app_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO your_app_user;

-- Log completion
SELECT 'Database initialization completed successfully' AS status;