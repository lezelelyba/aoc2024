resource "aws_vpc" "vpc" {
    cidr_block = var.vpc_cidr

    enable_dns_hostnames = true
    enable_dns_support = true

    tags = {
        name = "${var.env}-aoc2024-vpc"
        environment = var.env
    }
}

resource "aws_subnet" "public_subnet" {
    vpc_id = aws_vpc.vpc.id
    cidr_block = var.public_cidr
    map_public_ip_on_launch = true

    tags = {
        name = "${var.env}-aoc2024-public-subnet"
        environment = var.env
    }
}

resource "aws_subnet" "private_subnet" {
    vpc_id = aws_vpc.vpc.id
    cidr_block = var.private_cidr
    map_public_ip_on_launch = false 

    tags = {
        name = "${var.env}-aoc2024-private-subnet"
        environment = var.env
    }
}

resource "aws_internet_gateway" "igw" {
    vpc_id = aws_vpc.vpc.id

    tags = {
        name = "${var.env}-aoc2024-igw"
        environment = var.env
    }
}

resource "aws_eip" "nat_eip" {
    domain = "vpc" 
    depends_on = [aws_internet_gateway.igw]
}

resource "aws_nat_gateway" "nat" {
    allocation_id = aws_eip.nat_eip
    subnet_id = aws_subnet.private_subnet

    tags = {
        name = "${var.env}-aoc2024-natgw"
        environment = var.env
    }
}

resource "aws_route_table" "private" {
   vpc_id = aws_vpc.vpc.id 

   tags = {
       name = "${var.env}-aoc2024-private-rt"
       environment = var.env
   }
}

resource "aws_route" "private_gw" {
    route_table_id = aws_route_table.private.id
    destination_cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.nat.id
}

resource "aws_route_table" "public" {
   vpc_id = aws_vpc.vpc.id 
   
   tags = {
       name = "${var.env}-aoc2024-public-rt"
       environment = var.env
   }
}

resource "aws_route" "public_gw" {
    route_table_id = aws_route_table.public.id
    destination_cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw 
}

resource "aws_route_table_association" "private" {
   subnet_id = aws_subnet.private_subnet.id
   route_table_id = aws_route_table.private
}

resource "aws_route_table_association" "public" {
   subnet_id = aws_subnet.public_subnet.id
   route_table_id = aws_route_table.public
}