create table gym_consistency (
	Id int GENERATED ALWAYS AS identity PRIMARY KEY,
	date_gym varchar(10) unique,
	ok bool
);