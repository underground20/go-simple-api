db = db.getSiblingDB("admin");
db.createUser({
    user: process.env.DB_USER,
    pwd: process.env.DB_PASS,
    mechanisms: ["SCRAM-SHA-256"],
    roles: [
        { role: "readWrite", db: process.env.DB_NAME },
    ]
});

db = db.getSiblingDB(process.env.DB_NAME);
db.createCollection("employees");

print("Initialization completed: user and collections created.");