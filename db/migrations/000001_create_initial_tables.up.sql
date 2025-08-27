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
  payment_date TIMESTAMP DEFAULT NULL,
  invoice_id VARCHAR(255) NOT NULL,
  external_id VARCHAR(255) NOT NULL,
  invoice_url TEXT NOT NULL,
  status payment_status_enum NOT NULL DEFAULT 'PENDING',
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NULL,
  CONSTRAINT fk_payment_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id)
);
