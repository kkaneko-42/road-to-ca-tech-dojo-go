CREATE DATABASE IF NOT EXISTS road_to_ca DEFAULT CHARACTER SET utf8mb4;
USE road_to_ca;

CREATE TABLE `game_setting` (gacha_coin_consumption integer);
CREATE TABLE `users_infos` (user_id char(8) primary key, user_name varchar(16), having_coins integer);
CREATE TABLE `users_tokens` (user_id char(8) primary key, token char(16) unique);
CREATE TABLE `scores` (user_id char(8), score integer);
CREATE TABLE `users_inventories` (user_id char(8), item_id char(8));
CREATE TABLE `items` (item_id char(8) primary key, item_name varchar(32), rarity integer, gacha_weight integer);
