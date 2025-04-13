#!/usr/bin/env python3
import os
from pathlib import Path
import psycopg2
import random
from faker import Faker
from dotenv import load_dotenv

# Ids
from ulid import ULID
from ksuid import Ksuid
from uuid import uuid4
from nanoid import generate
from cuid2 import cuid
from snowflake import SnowflakeGenerator

# Initialize Faker
fake = Faker()

# Example departments list
departments = [
    "accounting", "sales", "engineering", "hr", "marketing", "it", "operations", "marketing", "operations"
]

def get_db_connection():
    # Get database connection parameters from environment variables.
    # Adjust defaults as needed.
    env_path = Path(__file__).resolve().parents[2] / '.env'
    print(env_path)
    # Load environment variables
    load_dotenv(dotenv_path=env_path)
    postgres_user = os.getenv("DATABASE_USERNAME", "postgres")
    postgres_db = os.getenv("DATABASE_NAME")
    postgres_password = os.getenv("DATABASE_PASSWORD")
    postgres_port = os.getenv("DATABASE_PORT")

    conn = psycopg2.connect(
        dbname=postgres_db,
        user=postgres_user,
        password=postgres_password,
        host="localhost",
        port=postgres_port
    )
    return conn

def create_users_table(conn, table_name):
    with conn.cursor() as cur:
        cur.execute(f"""
            CREATE {table_name} IF NOT EXISTS USERS (
                ID varchar(26) NOT NULL,
                HASH VARCHAR(64) NOT NULL,
                USER_NAME VARCHAR(50) NOT NULL,
                FIRST_NAME VARCHAR(255) NOT NULL,
                LAST_NAME VARCHAR(255) NOT NULL,
                EMAIL VARCHAR(255) NOT NULL,
                DEPARTMENT VARCHAR(255),
                PRIMARY KEY(ID),
                UNIQUE(EMAIL),
                UNIQUE(USER_NAME)
            );
        """)
        conn.commit()
        print("USERS table created (or already exists).")

def generate_fake_user(id):
    # Generate a random 64-character hex string (32 bytes * 2 hex digits each)
    hash_value = os.urandom(32).hex()
    user_name = fake.user_name()[:50]
    first_name = fake.first_name()[:255]
    last_name = fake.last_name()[:255]
    email = fake.email()[:255]
    # Example: random user status letter; adjust choices as needed.
    department = random.choice(departments)
    return (id, hash_value, user_name, first_name, last_name, email, department)

def insert_fake_users(conn, id, tables, num_records):
    with conn.cursor() as cur:
        for _ in range(num_records):
            user_data = generate_fake_user(id())
            cur.execute(f"""
                INSERT INTO {tables} (id, HASH, USER_NAME, FIRST_NAME, LAST_NAME, EMAIL, DEPARTMENT)
                VALUES (%s, %s, %s, %s, %s, %s, %s);
            """, user_data)
        conn.commit()
        print(f"Inserted {num_records} fake user records.")

def main():
    num_records = 1_000_000  # Adjust the number of records you want to genera
    conn = get_db_connection()
    try:
        ids = [(ULID, "UserUlid"), (uuid4, "UserUuid4"), (Ksuid, "UserKsuid"), (generate, "UserNanoid"), (cuid, "UserCuid"), (SnowflakeGenerator(), "UserSnowflake")]
        for id, table in ids:
            create_users_table(conn, table)
            insert_fake_users(conn, id, table, num_records)
    finally:
        conn.close()

if __name__ == '__main__':
    main()
