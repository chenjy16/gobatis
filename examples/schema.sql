-- 用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    real_name VARCHAR(100),
    age INT,
    gender ENUM('male', 'female', 'other') DEFAULT 'other',
    phone VARCHAR(20),
    department VARCHAR(100),
    position VARCHAR(100),
    salary DECIMAL(10,2),
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_department (department),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);

-- 部门表
CREATE TABLE departments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    manager_id BIGINT,
    budget DECIMAL(15,2),
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (manager_id) REFERENCES users(id)
);

-- 用户角色表
CREATE TABLE user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    role_name VARCHAR(50) NOT NULL,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    granted_by BIGINT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES users(id),
    UNIQUE KEY uk_user_role (user_id, role_name)
);

-- 用户登录日志表
CREATE TABLE user_login_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT,
    login_result ENUM('success', 'failed') DEFAULT 'success',
    failure_reason VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_login_time (login_time),
    INDEX idx_login_result (login_result)
);

-- 插入测试数据
INSERT INTO users (username, email, password, real_name, age, gender, phone, department, position, salary, status, last_login_at) VALUES
('admin', 'admin@example.com', 'hashed_password_1', '管理员', 35, 'male', '13800138000', 'IT', '系统管理员', 15000.00, 'active', '2024-01-15 10:30:00'),
('zhangsan', 'zhangsan@example.com', 'hashed_password_2', '张三', 28, 'male', '13800138001', 'IT', '高级开发工程师', 12000.00, 'active', '2024-01-15 09:15:00'),
('lisi', 'lisi@example.com', 'hashed_password_3', '李四', 32, 'female', '13800138002', 'HR', '人事经理', 10000.00, 'active', '2024-01-14 16:45:00'),
('wangwu', 'wangwu@example.com', 'hashed_password_4', '王五', 26, 'male', '13800138003', '研发部', '初级开发工程师', 8000.00, 'active', '2024-01-15 08:20:00'),
('zhaoliu', 'zhaoliu@example.com', 'hashed_password_5', '赵六', 29, 'female', '13800138004', '市场部', '市场专员', 7000.00, 'inactive', '2024-01-10 14:30:00'),
('sunqi', 'sunqi@example.com', 'hashed_password_6', '孙七', 31, 'male', '13800138005', 'IT', '架构师', 18000.00, 'active', '2024-01-15 11:00:00'),
('zhouba', 'zhouba@example.com', 'hashed_password_7', '周八', 24, 'female', '13800138006', 'HR', '招聘专员', 6000.00, 'active', '2024-01-13 13:15:00'),
('wujiu', 'wujiu@example.com', 'hashed_password_8', '吴九', 27, 'male', '13800138007', '财务部', '会计', 8500.00, 'active', '2024-01-12 17:20:00');

INSERT INTO departments (name, description, manager_id, budget, status) VALUES
('IT', 'Information Technology Department', 1, 500000.00, 'active'),
('HR', 'Human Resources Department', 3, 200000.00, 'active'),
('研发部', 'Research and Development Department', 6, 800000.00, 'active'),
('市场部', 'Marketing Department', NULL, 300000.00, 'active'),
('财务部', 'Finance Department', 8, 150000.00, 'active');

INSERT INTO user_roles (user_id, role_name, granted_by) VALUES
(1, 'ADMIN', 1),
(1, 'USER', 1),
(2, 'DEVELOPER', 1),
(2, 'USER', 1),
(3, 'HR_MANAGER', 1),
(3, 'USER', 1),
(4, 'DEVELOPER', 1),
(4, 'USER', 1),
(5, 'MARKETING', 1),
(5, 'USER', 1),
(6, 'ARCHITECT', 1),
(6, 'DEVELOPER', 1),
(6, 'USER', 1),
(7, 'HR', 1),
(7, 'USER', 1),
(8, 'FINANCE', 1),
(8, 'USER', 1);