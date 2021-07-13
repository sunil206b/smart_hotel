create table users(
    id serial primary key,
    first_name VARCHAR(255) not null,
    last_name VARCHAR(255) not null,
    email VARCHAR(255) not null unique,
    password text not null,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    access_level INTEGER
);

create INDEX idx_users_email ON users(email);

create table reservations(
    id serial primary key,
    first_name VARCHAR(255) not null,
    last_name VARCHAR(255) not null,
    email VARCHAR(255) not null,
    phone VARCHAR(20) not null,
    check_in DATE,
    check_out DATE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    room_id INT,
    foreign key(room_id) references rooms(id)
);

CREATE INDEX idx_reservations_email ON reservations(email);
CREATE INDEX idx_reservations_last_name ON reservations(last_name);

create table rooms(
    id serial primary key,
    room_name VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

create table room_restrictions(
    id serial primary key,
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    room_id INT,
    reservation_id INT,
    restriction_id INT,
    foreign key(reservation_id) references reservations(id),
    foreign key(room_id) references rooms(id),
    foreign key(restriction_id) references restrictions(id)
);

ALTER TABLE room_restrictions ALTER COLUMN reservation_id DROP NOT NULL;
create INDEX idx_room_restrictions_start_date ON room_restrictions(start_date);
create INDEX idx_room_restrictions_end_date ON room_restrictions(end_date);
create INDEX idx_room_restrictions_room_id ON room_restrictions(room_id);
create INDEX idx_room_restrictions_restriction_id ON room_restrictions(restriction_id);
create INDEX idx_room_restrictions_reservation_id ON room_restrictions(reservation_id);

create table restrictions(
    id serial primary key,
    restriction_name VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

select count(id) from room_restrictions where to_date('2021-07-13', 'YYYY-MM-DD') < date(end_date) and to_date('2021-07-18', 'YYYY-MM-DD') > date(start_date);

-- search date is exactly same as existing reservation
select count(id) from room_restrictions where '2021-07-13' < end_date and '2021-07-18' > start_date;

-- start date is before existing reservation, end date is same
select count(id) from room_restrictions where '2021-07-12' < end_date and '2021-07-18' > start_date;

-- end date is after existing reservation end date, start date is same
select count(id) from room_restrictions where '2021-07-13' < end_date and '2021-07-19' > start_date;

-- search dates are outside of all existing reservations, but cover the reservation
select count(id) from room_restrictions where '2021-07-10' < end_date and '2021-07-25' > start_date;

-- search dates are outside of all existing reservations
select count(id) from room_restrictions where '2021-07-19' < end_date and '2021-07-28' > start_date;

-- search dates are inside of existing reservations
select count(id) from room_restrictions where '2021-07-14' < end_date and '2021-07-17' > start_date;
