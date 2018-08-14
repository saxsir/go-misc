create table reports (
  id int(11) not null auto_increment,
  title varchar(255),
  body text,
  created datetime not null default now(),
  updated datetime ,
  primary key (`id`)
) ENGINE=InnoDB;
