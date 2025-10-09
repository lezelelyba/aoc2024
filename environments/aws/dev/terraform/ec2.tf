# data "aws_ami" "amazon_linux" {
#     most_recent = true
#     owners = ["amazon"]
# 
#     # name = "name" = filter which filters based on "name"
#     filter {
#         name = "name"
#         values = ["amzn2-ami-hvm-*-x86_64-gp2"]
#     }
# }
# 
# resource "aws_key_pair" "pub" {
#     key_name = "${var.env}-bastion-login"
#     public_key = file(var.sshpubkeypath)
# }
# 
# resource "aws_instance" "bastion" {
#     ami = data.aws_ami.amazon_linux.id
# 
#     instance_type = "t3.micro"
#     subnet_id = aws_subnet.public[0].id
# 
#     key_name = aws_key_pair.pub.key_name
# 
#     vpc_security_group_ids = [aws_security_group.bastion.id] 
# 
#     tags = {
#         name = "${var.env}-aoc2024-bastion-ec2"
#         environment = var.env
#     }
# }
# 
# resource "aws_instance" "test_host" {
#     ami = data.aws_ami.amazon_linux.id
# 
#     instance_type = "t3.micro"
#     subnet_id = aws_subnet.private[0].id
# 
#     key_name = aws_key_pair.pub.key_name
# 
#     vpc_security_group_ids = [aws_security_group.testhost.id] 
# 
#     tags = {
#         name = "${var.env}-aoc2024-testhost-ec2"
#         environment = var.env
#     }
# }