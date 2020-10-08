CREATE DATABASE IF NOT EXISTS testAPI;
USE testAPI;

CREATE TABLE IF NOT EXISTS songs(
	id TINYINT NOT NULL PRIMARY KEY AUTO_INCREMENT, 
    title VARCHAR(50) NOT NULL DEFAULT 'N/A',
    duration VARCHAR(50) NOT NULL DEFAULT 'N/A',
    singer VARCHAR(50) NOT NULL DEFAULT 'N/A'
);

# use  truncate to delete all records from a table
truncate table songs

select * from songs
