CREATE TABLE t_airports (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    iata_code VARCHAR(10),
    city VARCHAR(100),
    country_id INT,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_airlines (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    iata_code VARCHAR(10),
    logo_base64 TEXT
);

CREATE TABLE t_countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    iso_code VARCHAR(10),
    iso_code3 VARCHAR(10),
    currency_code VARCHAR(10),
    phone_code VARCHAR(10),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE t_backoffice_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE t_backoffice_users (
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    remember_token VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_contacts (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE t_flight_bookings (
    id SERIAL PRIMARY KEY,
    transaction_id INT,
    pnrid INT,
    booking_code VARCHAR(20),
    ntsa NUMERIC(15,2),
    commission NUMERIC(15,2),
    publish_fare NUMERIC(15,2),
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_flight_passengers (
    id SERIAL PRIMARY KEY,
    booking_id INT,
    title VARCHAR(10),
    full_name VARCHAR(255),
    type VARCHAR(20),
    pax_id INT,

    dob DATE,
    adult_assoc INT,

    passport_no VARCHAR(50),
    passport_nationality VARCHAR(100),
    passport_issuing_country VARCHAR(100),
    passport_expiry DATE,

    ticket_number VARCHAR(50),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_flight_remarks (
    id SERIAL PRIMARY KEY,
    booking_id INT NOT NULL,
    remark TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_flight_segments (
    id SERIAL PRIMARY KEY,
    booking_id INT,
    airlinecode VARCHAR(10),
    flightno VARCHAR(255),
    fromcity VARCHAR(255),
    tocity VARCHAR(255),
    classcode VARCHAR(10),
    cabin VARCHAR(255),
    departtime TIME,
    arrivetime TIME,
    departdate VARCHAR(255),
    arrivaldate VARCHAR(255),
    durationhour INT,
    durationminute INT,
    FOREIGN KEY (booking_id) REFERENCES t_flight_bookings(id)
);

CREATE TABLE t_flight_transactions (
    id SERIAL PRIMARY KEY,
    contact_id INT,
    vm_code VARCHAR(20),
    payment_method_id INT,

    amount NUMERIC(15,2),
    admin_fee NUMERIC(15,2),
    markup NUMERIC(15,2),
    payment_gateway_charge NUMERIC(15,2),

    trip_type VARCHAR(20),
    status VARCHAR(20),
    source_platform VARCHAR(20),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_payment_methods (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_phone_balance_transactions (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    type VARCHAR(20) NOT NULL,
    processed_by INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_va_banks (
    id SERIAL PRIMARY KEY,
    payment_method_id INT NOT NULL,
    bank_name VARCHAR(100) NOT NULL,
    bank_code VARCHAR(50) NOT NULL,
    logo_base64 TEXT,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_vending_machine_payment_methods (
    id SERIAL PRIMARY KEY,
    vending_machine_id INT NOT NULL,
    payment_method_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE t_vending_machines (
    id SERIAL PRIMARY KEY,
    vm_code VARCHAR(50) NOT NULL,
    airport_id INT NOT NULL,

    admin_fee NUMERIC(15,2) DEFAULT 0,
    markup_pesawat NUMERIC(15,2) DEFAULT 0,
    markup_kereta NUMERIC(15,2) DEFAULT 0,
    markup_hotel NUMERIC(15,2) DEFAULT 0,
    markup_ppob NUMERIC(15,2) DEFAULT 0,
    markup_themepark NUMERIC(15,2) DEFAULT 0,

    is_active BOOLEAN DEFAULT TRUE,
    last_heartbeat TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    username VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    full_name VARCHAR(255),
    phone VARCHAR(255),
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);