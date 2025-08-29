#!/usr/bin/env python3
import os
import asyncio
from pathlib import Path
import random
import time
from faker import Faker
from dotenv import load_dotenv

# Id generators
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

# Create generators if needed
cuid_generator = Cuid(length=25)
snow_flake_generator = SnowflakeGenerator(25)

async def get_db_connection():
    # Determine the .env file location relative to this script.
    env_path = Path(__file__).resolve().parents[2] / '.env'
    print("Loading env from:", env_path)
    load_dotenv(dotenv_path=env_path)

    postgres_user = os.getenv("DATABASE_USERNAME", "postgres")
    postgres_db = os.getenv("DATABASE_NAME")
    postgres_password = os.getenv("DATABASE_PASSWORD")
    postgres_port = os.getenv("DATABASE_PORT")
    postgres_host = "localhost"

    # Construct connection string using URL (URI) format.
    conn_str = (
        f"postgresql://{postgres_user}:{postgres_password}"
        f"@{postgres_host}:{postgres_port}/{postgres_db}?sslmode=disable"
    )

    # Use asyncpg to connect asynchronously.
    import asyncpg  # Import here if not installed globally.
    conn = await asyncpg.connect(conn_str)
    return conn


def generate_unique_email(unique_id, first_name):
    # Use Faker for a realistic name and append the unique ID to ensure uniqueness.
    name = first_name.lower()
    domain = fake.free_email_domain()
    return f"{name}{unique_id}@{domain}"

def generate_fake_user(id):
    # Generate fake user data using Faker.
    user_name = fake.user_name()[:50]
    first_name = fake.first_name()[:255]
    last_name = fake.last_name()[:255]
    email = generate_unique_email(id, first_name)
    department = random.choice(departments)
    # Return a tuple corresponding to the table columns.
    return (id, user_name, first_name, last_name, email, department)

async def insert_fake_users(conn, table, num_records):
    # Insert 'num_records' rows into the given table.
    for _ in range(num_records):
        user_data = None
        if table == "users_ulid":
            # ULID expects an integer timestamp in milliseconds.
            ulid = ULID()
            id = str(ulid)
            print(f"Generated ULID: {str(id)} length: {len(id)}")
            user_data = generate_fake_user(id)
        elif table == "users_uuid":
            id = uuid4()
            user_data = generate_fake_user(str(id))
        elif table == "users_ksuid":
            id = Ksuid()
            user_data = generate_fake_user(str(id))
        elif table == "users_nanoid":
            id = generate(size=21)
            user_data = generate_fake_user(id)
        elif table == "users_cuid":
            id = cuid_generator.generate()
            user_data = generate_fake_user(id)
        elif table == "users_snowflake":
            id = next(snow_flake_generator)
            user_data = generate_fake_user(id)
        # Execute the insert asynchronously.
        await conn.execute(
            f"""
            INSERT INTO {table}
              (id, USER_NAME, FIRST_NAME, LAST_NAME, EMAIL, DEPARTMENT)
            VALUES ($1, $2, $3, $4, $5, $6)
            """,
            user_data[0],
            user_data[1],
            user_data[2],
            user_data[3],
            user_data[4],
            user_data[5]
        )
    print(f"Inserted {num_records} fake user records into table {table}.")

async def main():
    num_records = 1_000_000  # Adjust the number of records as needed.
    conn = await get_db_connection()
    try:
        tables = [
            "users_ulid",
            "users_uuid",
            "users_ksuid",
            "users_nanoid",
            "users_cuid",
            "users_snowflake"
        ]
        # Loop through each table and insert fake users asynchronously.
        for table in tables:
            await insert_fake_users(conn, table, num_records)
    finally:
        await conn.close()

if __name__ == '__main__':
    asyncio.run(main())
