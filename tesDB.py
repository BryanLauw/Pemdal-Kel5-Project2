import os

# ===================================
# 🔧 Cassandra Driver Safe Mode (Python 3.12)
# ===================================
os.environ["CASS_DRIVER_NO_EXTENSIONS"] = "1"
os.environ["CASS_DRIVER_NO_CYTHON"] = "1"

from cassandra.cluster import Cluster
from cassandra.connection import DefaultConnection
from cassandra.policies import RoundRobinPolicy
from neo4j import GraphDatabase

# Force Cassandra to use pure-Python socket (tanpa asyncore/libev)
Cluster.connection_class = DefaultConnection


# ===================================
# 🔹 CASSANDRA CRUD TESTING
# ===================================
def test_cassandra():
    print("\n=== 🧩 Testing Cassandra CRUD ===")

    # Ganti "localhost" jadi "cassandra" jika kamu pakai docker-compose
    cluster = Cluster(["localhost"], port=9042, load_balancing_policy=RoundRobinPolicy())
    session = cluster.connect()

    # Buat keyspace
    session.execute("""
        CREATE KEYSPACE IF NOT EXISTS test_app
        WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
    """)

    # Gunakan keyspace
    session.set_keyspace("test_app")

    # Buat tabel
    session.execute("""
        CREATE TABLE IF NOT EXISTS users (
            email TEXT PRIMARY KEY,
            name TEXT,
            city TEXT
        );
    """)

    # CREATE
    session.execute("""
        INSERT INTO users (email, name, city)
        VALUES (%s, %s, %s)
    """, ("budi@example.com", "Budi", "Bandung"))
    print("✅ Data inserted")

    # READ
    rows = session.execute("SELECT * FROM users;")
    for row in rows:
        print("📄 Read:", row.email, row.name, row.city)

    # UPDATE
    session.execute("""
        UPDATE users SET city = %s WHERE email = %s
    """, ("Jakarta", "budi@example.com"))
    print("✏️ Data updated")

    # READ AGAIN
    row = session.execute("SELECT * FROM users WHERE email=%s", ("budi@example.com",)).one()
    print("📄 After update:", row.email, row.name, row.city)

    # DELETE
    session.execute("DELETE FROM users WHERE email=%s", ("budi@example.com",))
    print("🗑️ Data deleted")

    # VERIFY DELETE
    rows = session.execute("SELECT * FROM users;")
    print("Remaining rows:", list(rows))

    cluster.shutdown()
    print("✅ Cassandra test done.")


# ===================================
# 🔹 NEO4J CRUD TESTING
# ===================================
def test_neo4j():
    print("\n=== 🧩 Testing Neo4j CRUD ===")

    # Ganti host sesuai docker-compose (biasanya "neo4j")
    uri = "bolt://localhost:7687"
    driver = GraphDatabase.driver(uri, auth=("neo4j", "password123"))

    with driver.session(database="neo4j") as session:
        # CREATE node
        session.run("""
            CREATE (u:User {email: $email, name: $name, city: $city})
        """, {"email": "andi@example.com", "name": "Andi", "city": "Surabaya"})
        print("✅ Node created")

        # READ node
        result = session.run("MATCH (u:User) RETURN u.email AS email, u.name AS name, u.city AS city")
        for record in result:
            print("📄 Read:", record["email"], record["name"], record["city"])

        # UPDATE node
        session.run("""
            MATCH (u:User {email: $email})
            SET u.city = $city
        """, {"email": "andi@example.com", "city": "Yogyakarta"})
        print("✏️ Node updated")

        # READ updated node
        record = session.run("""
            MATCH (u:User {email: $email}) RETURN u
        """, {"email": "andi@example.com"}).single()
        print("📄 After update:", record["u"])

        # DELETE node
        session.run("MATCH (u:User {email: $email}) DELETE u", {"email": "andi@example.com"})
        print("🗑️ Node deleted")

        # VERIFY delete
        result = session.run("MATCH (u:User) RETURN count(u) AS count").single()
        print("Remaining nodes:", result["count"])

    driver.close()
    print("✅ Neo4j test done.")


# ===================================
# 🚀 MAIN EXECUTION
# ===================================
if __name__ == "__main__":
    test_cassandra()
    test_neo4j()
