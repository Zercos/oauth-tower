create table IF NOT EXISTS clients (client_id text not null primary key, client_secret text, redirect_uri text);
create table IF NOT EXISTS users (user_id integer primary key, username text, password_hash text);

INSERT INTO clients (client_id, client_secret, redirect_uri)
SELECT 'test_client', 'secret', 'http://localhost:3000/callback'
WHERE NOT EXISTS(SELECT 1 FROM clients WHERE client_id = 'test_client');

INSERT INTO users (username, password_hash)
SELECT 'test_user', '33'
WHERE NOT EXISTS(SELECT 1 FROM users WHERE username = 'test_user');