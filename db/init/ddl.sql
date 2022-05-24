CREATE DATABASE IF NOT EXISTS road_to_ca;

USE road_to_ca;

CREATE TABLE `game_setting` (gacha_coin_consumption integer);
CREATE TABLE `users_infos` (user_id char(4), user_name varchar(16), having_coins integer);
CREATE TABLE `users_tokens` (user_id char(4), token varchar(32));
CREATE TABLE `scores` (user_id char(4), score integer);
CREATE TABLE `users_inventories` (user_id char(4), item_id char(4));
CREATE TABLE `items` (item_id char(4), item_name varchar(16), rarity integer, gacha_weight integer);
