// Create or switch to a new database
db = db.getSiblingDB('my_database');

// Create a new collection
db.createCollection('users');

// Insert a document into the 'users' collection
db.users.insertOne({
    username: "johndoe",
    email: "johndoe@example.com",
    password_hash: "$2a$10$eImiTXuWVxfM37uY4JANjQeG9R5E1W2z3obEW6TlvEc4cGnOKH4WC", // Example bcrypt hash
    is_admin: true,
    created_at: new Date(),
    updated_at: new Date()
});