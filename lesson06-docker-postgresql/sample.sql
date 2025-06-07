-- Dữ liệu mẫu

-- =======================================================
-- Phần 1: Mối quan hệ Một - Một (One-to-One)
-- =======================================================

-- Chèn dữ liệu vào bảng 'users' trước
INSERT INTO users (name, email) VALUES
('Vũ Quốc Tuấn', 'tuan@example.com'),
('Trần Thị Bình', 'binh.tran@example.com'),
('Lê Văn Cường', 'cuong.le@example.com'),
('Phạm Vân', 'van.pham@example.com');

-- Chèn dữ liệu vào bảng 'profiles', liên kết với 'users'
-- Lưu ý: user_id (1, 2, 3, 4) tương ứng với các user vừa được tạo ở trên
INSERT INTO profiles (user_id, phone, address) VALUES
(1, '0987654321', '123 Đường Láng, Đống Đa, Hà Nội'),
(2, '0912345678', '456 Lê Lợi, Quận 1, TP. Hồ Chí Minh'),
(3, '0334455667', '789 Ngô Quyền, Sơn Trà, Đà Nẵng'),
(4, '0888999777', '101 Hùng Vương, Ninh Kiều, Cần Thơ');


-- =======================================================
-- Phần 2: Mối quan hệ Một - Nhiều (One-to-Many)
-- =======================================================

-- Chèn dữ liệu vào bảng 'categories' trước
INSERT INTO categories (name) VALUES
('Đồ điện tử'),
('Quần áo'),
('Sách'),
('Đồ gia dụng');

-- Chèn dữ liệu vào bảng 'products', liên kết với 'categories'
-- category_id: 1='Đồ điện tử', 2='Quần áo', 3='Sách', 4='Đồ gia dụng'
-- status: 1=Còn hàng, 2=Hết hàng
INSERT INTO products (category_id, name, price, image, status) VALUES
(1, 'Điện thoại thông minh Galaxy S24', 21000000, 'images/galaxy_s24.jpg', 1),
(1, 'Laptop Dell XPS 15', 35000000, 'images/dell_xps15.jpg', 1),
(1, 'Tai nghe không dây Sony WH-1000XM5', 7500000, 'images/sony_xm5.jpg', 2),
(2, 'Áo thun nam cổ tròn', 250000, 'images/tshirt_basic.jpg', 1),
(2, 'Quần Jeans nữ ống rộng', 550000, 'images/jeans_wide.jpg', 1),
(3, 'Tiểu thuyết "Nhà Giả Kim"', 85000, 'images/nhagiakim.jpg', 1),
(3, 'Sách "Lược sử loài người"', 150000, 'images/luocsu.jpg', 1),
(4, 'Nồi chiên không dầu Philips', 3200000, 'images/philips_airfryer.jpg', 1);


-- =======================================================
-- Phần 3: Mối quan hệ Nhiều - Nhiều (Many-to-Many)
-- =======================================================

-- Chèn dữ liệu vào bảng 'students'
INSERT INTO students (name) VALUES
('Vũ Quốc Tuấn'),
('Bùi Bích Phương'),
('Đặng Quốc Bảo'),
('Vũ Minh Thư');

-- Chèn dữ liệu vào bảng 'courses'
INSERT INTO courses (name) VALUES
('Cơ sở dữ liệu'),
('Lập trình Web'),
('Trí tuệ nhân tạo'),
('Mạng máy tính');

-- Chèn dữ liệu vào bảng liên kết 'students_courses'
-- Tạo mối quan hệ giữa sinh viên và khóa học
INSERT INTO students_courses (student_id, course_id) VALUES
-- Tuấn (id=1) học khóa Cơ sở dữ liệu (id=1) và Lập trình Web (id=2)
(1, 1),
(1, 2),

-- Hà (id=2) học khóa Lập trình Web (id=2) và Trí tuệ nhân tạo (id=3)
(2, 2),
(2, 3),

-- Bảo (id=3) chỉ học khóa Cơ sở dữ liệu (id=1)
(3, 1),

-- Linh (id=4) học tất cả các khóa
(4, 1),
(4, 2),
(4, 3),
(4, 4);