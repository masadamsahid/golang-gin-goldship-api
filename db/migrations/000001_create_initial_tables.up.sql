CREATE TYPE role_enum AS ENUM(
  'SUPERADMIN',
  'ADMIN',
  'EMPLOYEE',
  'COURIER',
  'CUSTOMER'
);

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role role_enum NOT NULL DEFAULT 'CUSTOMER',
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS branches (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  phone TEXT NOT NULL,
  address TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NULL
);

CREATE TYPE shipment_status_enum AS ENUM(
  'PENDING_PAYMENT',
  'READY_TO_PICKUP',
  'PICKED_UP',
  'IN_TRANSIT',
  'DELIVERED',
  'CANCELLED'
);

CREATE TABLE IF NOT EXISTS shipments (
  id SERIAL PRIMARY KEY,
  tracking_number VARCHAR(255) UNIQUE,
  sender_id INT NOT NULL,
  sender_name VARCHAR(255) NOT NULL,
  sender_phone VARCHAR(20) NOT NULL,
  sender_address TEXT NOT NULL,
  recipient_name VARCHAR(255) NOT NULL,
  recipient_address TEXT NOT NULL,
  recipient_phone VARCHAR(20) NOT NULL,
  item_name VARCHAR(255) NOT NULL,
  item_weight DECIMAL(10, 2) NOT NULL,
  distance DECIMAL(10, 2) NOT NULL,
  status shipment_status_enum NOT NULL DEFAULT 'PENDING_PAYMENT',
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_shipments_user FOREIGN KEY (sender_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS shipment_histories (
  id SERIAL PRIMARY KEY,
  shipment_id INT NOT NULL,
  status shipment_status_enum NOT NULL,
  "desc" TEXT,
  courier_id INT,
  branch_id INT,
  timestamp TIMESTAMP DEFAULT NOW() NOT NULL,
  CONSTRAINT fk_shipment_histories_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id),
  CONSTRAINT fk_shipment_histories_courier FOREIGN KEY (courier_id) REFERENCES users(id),
  CONSTRAINT fk_shipment_histories_branch FOREIGN KEY (branch_id) REFERENCES branches(id)
);

CREATE TYPE payment_status_enum AS ENUM (
  'PENDING',
  'PAID',
  'EXPIRED',
  'CANCELLED'
);

CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  shipment_id INT NOT NULL,
  amount INT NOT NULL,
  invoice_id VARCHAR(255) UNIQUE NOT NULL,
  external_id VARCHAR(255) UNIQUE NOT NULL,
  invoice_url TEXT UNIQUE NOT NULL,
  status payment_status_enum NOT NULL DEFAULT 'PENDING',
  paid_at TIMESTAMP DEFAULT NULL,
  expired_at TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_payment_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id)
);


-- SEEDER
INSERT INTO branches (name, phone, address)
VALUES ('Jakarta', '0211234567', 'Jl. Sudirman No. 1'),
('Surabaya', '0311234567', 'Jl. Thamrin No. 1'),
('Bandung', '0221234567', 'Jl. Asia Afrika No. 1'),
('Semarang', '0241234567', 'Jl. Pahlawan No. 1'),
('Yogyakarta', '0274123456', 'Jl. Malioboro No. 1'),
('Malang', '0341123456', 'Jl. Ijen No. 1'),
('Solo', '0271123456', 'Jl. Slamet Riyadi No. 1'),
('Cirebon', '0231123456', 'Jl. Kartini No. 1'),
('Bogor', '0251123456', 'Jl. Pajajaran No. 1'),
('Depok', '0217654321', 'Jl. Margonda Raya No. 1'),
('Tangerang', '0212345678', 'Jl. Raya Serpong No. 1'),
('Bekasi', '0218765432', 'Jl. Ahmad Yani No. 1'),
('Kediri', '0354123456', 'Jl. Dhoho No. 1'),
('Madiun', '0351123456', 'Jl. Pahlawan No. 1'),
('Purwokerto', '0281123456', 'Jl. Jend. Sudirman No. 1'),
('Pekalongan', '0285123456', 'Jl. Hayam Wuruk No. 1'),
('Tegal', '0283123456', 'Jl. Gajah Mada No. 1'),
('Salatiga', '0298123456', 'Jl. Diponegoro No. 1'),
('Magelang', '0293123456', 'Jl. Pahlawan No. 1'),
('Sukabumi', '0266123456', 'Jl. Ahmad Yani No. 1'),
('Tasikmalaya', '0265123456', 'Jl. HZ. Mustofa No. 1'),
('Jember', '0331123456', 'Jl. Kalimantan No. 1'),
('Banyuwangi', '0333123456', 'Jl. Gatot Subroto No. 1'),
('Lumajang', '0334123456', 'Jl. PB. Sudirman No. 1'),
('Probolinggo', '0335123456', 'Jl. Soekarno Hatta No. 1'),
('Pasuruan', '0343123456', 'Jl. Pahlawan No. 1'),
('Mojokerto', '0321123456', 'Jl. Majapahit No. 1'),
('Gresik', '0313951234', 'Jl. Dr. Wahidin Sudirohusodo No. 1'),
('Denpasar', '0361123456', 'Jl. Teuku Umar No. 1'),
('Singaraja', '0362123456', 'Jl. Ahmad Yani No. 1'),
('Ubud', '0361975678', 'Jl. Raya Ubud No. 1');

