-- phpMyAdmin SQL Dump
-- version 4.5.1
-- http://www.phpmyadmin.net
--
-- Host: 127.0.0.1
-- Generation Time: 2016-11-22 09:29:59
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
  `parent` bigint(20) NOT NULL,
  `star` int(2) NOT NULL,
  `sha1` varchar(100) COLLATE utf8_bin NOT NULL,
  `src` varchar(600) COLLATE utf8_bin NOT NULL,
  `source` varchar(300) COLLATE utf8_bin NOT NULL,
  `url` varchar(300) COLLATE utf8_bin NOT NULL,
  `name` varchar(600) COLLATE utf8_bin NOT NULL,
  `file_type` varchar(300) COLLATE utf8_bin NOT NULL,
  `size` int(11) NOT NULL,
  `coll_time` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `user`
--

CREATE TABLE `user` (
  `id` int(10) UNSIGNED NOT NULL,
  `username` varchar(300) COLLATE utf8_bin NOT NULL,
  `password` varchar(300) COLLATE utf8_bin NOT NULL,
  `last_ip` varchar(300) COLLATE utf8_bin NOT NULL,
  `last_time` datetime NOT NULL,
  `is_disabled` int(1) COLLATE utf8_bin NOT NULL
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
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `coll`
--
ALTER TABLE `coll`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;
--
-- 使用表AUTO_INCREMENT `user`
--
ALTER TABLE `user`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;


--
-- 转存表中的数据 `user`
--

INSERT INTO `user` (`id`, `username`, `password`, `last_ip`, `last_time`, `is_disabled`) VALUES
(1, 'admin@admin.com', 'dd94709528bb1c83d08f3088d4043f4742891f4f', '218.26.1.182', '2016-11-29 15:57:56', 0);