// used image
data "aws_ami" "amazon_linux" {
    most_recent = true
    owners = ["amazon"]

    # name = "name" = filter which filters based on "name"
    filter {
        name = "name"
        values = ["amzn2-ami-hvm-*-x86_64-gp2"]
    }
}

// define new key pair to use for login to bastion and private host
resource "aws_key_pair" "pub" {
    count = local.keys_required ? 1 : 0

    key_name = "${var.env}-aoc-ec2-login"
    public_key = file(var.sshpubkeypath)
}

// spin up bastion
resource "aws_instance" "bastion" {
    count = var.bastion ? 1 : 0

    ami = data.aws_ami.amazon_linux.id

    instance_type = "t3.micro"
    subnet_id = aws_subnet.public[0].id

    key_name = aws_key_pair.pub[0].key_name

    vpc_security_group_ids = [aws_security_group.bastion.id] 

    tags = {
        name = "${var.env}-aoc2024-bastion-ec2"
        environment = var.env
    }
}

// spin up private host
resource "aws_instance" "private_host" {
    count = var.private_host ? 1 : 0

    ami = data.aws_ami.amazon_linux.id

    instance_type = "t3.micro"
    subnet_id = aws_subnet.private[0].id

    key_name = aws_key_pair.pub[0].key_name

    vpc_security_group_ids = [aws_security_group.testhost.id] 

    tags = {
        name = "${var.env}-aoc2024-private-host-ec2"
        environment = var.env
    }
}