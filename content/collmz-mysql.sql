-- phpMyAdmin SQL Dump
-- version 4.5.1
-- http://www.phpmyadmin.net
--
-- Host: 127.0.0.1
-- Generation Time: 2016-11-20 15:04:49
-- 服务器版本： 10.1.16-MariaDB
-- PHP Version: 5.6.24

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `collmz`
--

-- --------------------------------------------------------

--
-- 表的结构 `coll`
--

CREATE TABLE `coll` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `sha1` varchar(100) COLLATE utf8_bin NOT NULL,
  `src` varchar(600) COLLATE utf8_bin NOT NULL,
  `source` varchar(300) COLLATE utf8_bin NOT NULL,
  `url` varchar(300) COLLATE utf8_bin NOT NULL,
  `name` varchar(600) COLLATE utf8_bin NOT NULL,
  `type` varchar(300) COLLATE utf8_bin NOT NULL,
  `size` int(11) NOT NULL,
  `coll_time` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `coll`
--
ALTER TABLE `coll`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `coll`
--
ALTER TABLE `coll`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
