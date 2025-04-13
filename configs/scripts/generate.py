#!/usr/bin/env python3
import os
from pathlib import Path
import psycopg2
import random
from faker import Faker
from dotenv import load_dotenv
import time, datetime

# Ids
from ulid import ULID
from ksuid import Ksuid
from uuid import uuid4
from nanoid import generate
from cuid2 import Cuid
from snowflake import SnowflakeGenerator

# Initialize Faker
fake = Faker()

# Example departments list
departments = [
    "accounting", "sales", "engineering", "hr", "marketing", "it", "operations", "marketing", "operations"
]

cuid_generator = Cuid(length=25)
snow_flake_generator = SnowflakeGenerator(25)

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

def generate_fake_user(id):
    # Generate a random 64-character hex string (32 bytes * 2 hex digits each)
    user_name = fake.user_name()[:50]
    first_name = fake.first_name()[:255]
    last_name = fake.last_name()[:255]
    email = fake.email()[:255]
    # Example: random user status letter; adjust choices as needed.
    department = random.choice(departments)
    return (id, user_name, first_name, last_name, email, department)

def insert_fake_users(conn, tables, num_records):
    with conn.cursor() as cur:
        for _ in range(num_records):
            user_data = None
            if tables == "users_ulid":
                id = ULID().encode_timestamp(int(time.time() * 1000))
                user_data = generate_fake_user(str(id))
            elif tables == "users_uuid4":
                id = uuid4()
                user_data = generate_fake_user(str(id))
            elif tables == "users_ksuid":
                id = Ksuid()
                user_data = generate_fake_user(str(id))
            elif tables == "users_nanoid":
                id = generate(size=21)
                user_data = generate_fake_user(id)
            elif tables == "users_cuid":
                id = cuid_generator.generate()
                user_data = generate_fake_user(id)
            elif tables == "users_snowflake":
                id = next(snow_flake_generator)
                user_data = generate_fake_user(id)
            print (user_data)
            # cur.execute(f"""
            #     INSERT INTO {tables} (id, HASH, USER_NAME, FIRST_NAME, LAST_NAME, EMAIL, DEPARTMENT)
            #     VALUES (%s, %s, %s, %s, %s, %s, %s);
            # """, user_data)
        # conn.commit()
        print(f"Inserted {num_records} fake user records.")

def main():
    num_records = 10  # Adjust the number of records you want to genera
    conn = get_db_connection()
    try:
        tables = ["users_ulid", "users_uuid4", "users_ksuid", "users_nanoid", "users_cuid", "users_snowflake"]
        for table in tables:
            insert_fake_users(conn, table, num_records)


    finally:
        conn.close()

if __name__ == '__main__':
    main()
