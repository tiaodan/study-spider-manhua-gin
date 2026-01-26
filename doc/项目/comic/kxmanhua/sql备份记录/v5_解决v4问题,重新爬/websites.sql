/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80407
 Source Host           : localhost:3306
 Source Schema         : comic

 Target Server Type    : MySQL
 Target Server Version : 80407
 File Encoding         : 65001

 Date: 14/01/2026 18:14:16
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for websites
-- ----------------------------
DROP TABLE IF EXISTS `websites`;
CREATE TABLE `websites`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `domain` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `need_proxy` tinyint(1) NOT NULL,
  `is_https` tinyint(1) NOT NULL,
  `is_refer` tinyint(1) NOT NULL,
  `cover_url_is_need_https` tinyint(1) NOT NULL,
  `chapter_content_url_is_need_https` tinyint(1) NOT NULL,
  `cover_url_concat_rule` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `chapter_content_url_concat_rule` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `cover_domain` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `chapter_content_domain` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `book_can_spider_type` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `chapter_can_spider_type` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `book_spider_req_body_eg_server_filepath` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `chapter_spider_req_body_eg_server_filepath` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `star_type` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `website_type_id` bigint NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_website_unique`(`name` ASC, `domain` ASC) USING BTREE,
  INDEX `fk_websites_website_type`(`website_type_id` ASC) USING BTREE,
  CONSTRAINT `fk_websites_website_type` FOREIGN KEY (`website_type_id`) REFERENCES `website_types` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `chk_websites_book_can_spider_type` CHECK (`book_can_spider_type` <> _utf8mb4''),
  CONSTRAINT `chk_websites_chapter_can_spider_type` CHECK (`chapter_can_spider_type` <> _utf8mb4''),
  CONSTRAINT `chk_websites_chapter_content_domain` CHECK (`chapter_content_domain` <> _utf8mb4''),
  CONSTRAINT `chk_websites_chapter_content_url_concat_rule` CHECK (`chapter_content_url_concat_rule` <> _utf8mb4''),
  CONSTRAINT `chk_websites_cover_domain` CHECK (`cover_domain` <> _utf8mb4''),
  CONSTRAINT `chk_websites_cover_url_concat_rule` CHECK (`cover_url_concat_rule` <> _utf8mb4''),
  CONSTRAINT `chk_websites_domain` CHECK (`domain` <> _utf8mb4''),
  CONSTRAINT `chk_websites_name` CHECK (`name` <> _utf8mb4'')
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of websites
-- ----------------------------
INSERT INTO `websites` VALUES (1, '待分类', '未知', 0, 0, 0, 1, 1, '{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}', '{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}', 'www.未知.com', 'www.未知.com', 'both', 'both', '爬json:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byhtml.html', '爬json:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byhtml.html', 'copy', 1);
INSERT INTO `websites` VALUES (2, 'j88d', 'www.j88d.com', 0, 0, 1, 0, 0, '{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}', '{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}', 'www.j88d.com', 'www.j88d.com', 'both', 'both', '爬json:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byhtml.html', '爬json:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byhtml.html', 'copy-toptoon', 8);
INSERT INTO `websites` VALUES (3, 'aws-s3', 'ap-northeast-2.console.aws.amazon.com/s3/home?region=ap-northeast-2', 0, 1, 1, 0, 0, '{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}', '{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}', 'www.awsS3.com', 'www.awsS3.com', 'json', 'json', '爬json:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byhtml.html', '爬json:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byhtml.html', 'my', 7);
INSERT INTO `websites` VALUES (4, '预留', 'www.yuliu.com', 0, 0, 0, 0, 0, '{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}', '{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}', 'www.预留.com', 'www.预留.com', 'html', 'html', '爬json:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byhtml.html', '爬json:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byhtml.html', 'my', 1);
INSERT INTO `websites` VALUES (5, '开心看漫画', 'kxmanhua.com', 0, 1, 0, 0, 0, '{website表-protocol}://{website表-domain}/{book表-cover_url_api_path}', '{website表-protocol}://{website表-domain}/{book表-chapter_content_url_api_path}', 'www.预留.com', 'www.预留.com', 'html', 'html', '爬json:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/book_spider_req_body_eg_byhtml.html', '爬json:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byjson.json; 爬html:doc/项目/comic/toptoon-tw/chapter_spider_req_body_eg_byhtml.html', 'my', 1);

SET FOREIGN_KEY_CHECKS = 1;
