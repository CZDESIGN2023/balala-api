ALTER TABLE `third_pf_account`
    ADD COLUMN `pf_name` varchar(255) NULL COMMENT '平台名' AFTER `pf_code`;


update third_pf_account set pf_name = 'IMChat' where pf_code = 3;
update third_pf_account set pf_name = '轻聊' where pf_code = 4;
update third_pf_account set pf_name = 'Halala' where pf_code = 5;