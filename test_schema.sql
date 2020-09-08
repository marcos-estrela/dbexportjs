CREATE TABLE test_table (
  id INT NOT NULL,
  name varchar(100) NULL
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE TABLE users (
  id INT NOT NULL,
  username VARCHAR(100) NOT NULL
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE OR REPLACE VIEW vw_users AS
SELECT * FROM users;

CREATE OR REPLACE VIEW vw_test_table AS
SELECT t.id, t.name, u.username FROM test_table t
JOIN users u ON u.id = t.id;

CREATE PROCEDURE `sp_users_ins`(IN p_username VARCHAR(100), INOUT p_status BOOLEAN, INOUT p_msg TEXT)
BEGIN
  DECLARE v_last_insert_id INT(11);

  INSERT INTO users (username) VALUES (p_username);

  SET v_last_insert_id = LAST_INSERT_ID();

  SET p_status = TRUE;
  SET p_msg = 'User registered successfully.';
END;
