// allow all available zones to be used
data "aws_availability_zones" "available" {
    state = "available"
}

// create vpc
resource "aws_vpc" "vpc" {
    cidr_block = var.vpc_cidr

    enable_dns_hostnames = true
    enable_dns_support = true

    tags = {
        name = "${var.env}-aoc2024-vpc"
        environment = var.env
    }
}

// create 2 public subnets
resource "aws_subnet" "public" {
    count = 2

    vpc_id = aws_vpc.vpc.id
    availability_zone = data.aws_availability_zones.available.names[count.index]
    // +4 subnet mask, 0000
    cidr_block = cidrsubnet(var.vpc_cidr, 4, count.index)
    map_public_ip_on_launch = true

    tags = {
        name = "${var.env}-aoc2024-public-subnet-${count.index}"
        environment = var.env
    }
}

// create 2 private subnets
resource "aws_subnet" "private" {
    count = 2

    vpc_id = aws_vpc.vpc.id
    availability_zone = data.aws_availability_zones.available.names[count.index]
    // +4 subnet mask, 1000
    cidr_block = cidrsubnet(var.vpc_cidr, 4, count.index + 8)
    map_public_ip_on_launch = false 

    tags = {
        name = "${var.env}-aoc2024-private-subnet-${count.index}"
        environment = var.env
    }
}

// create internet gateway for vpc
resource "aws_internet_gateway" "igw" {
    vpc_id = aws_vpc.vpc.id

    tags = {
        name = "${var.env}-aoc2024-igw"
        environment = var.env
    }
}

// create elastic ip for nat gateways
resource "aws_eip" "nat_eip" {
    count = 2

    domain = "vpc" 
    depends_on = [aws_internet_gateway.igw]
    tags = {
        name = "${var.env}-aoc2024-eip-${count.index}"
        environment = var.env
    }
}

// create nat gateways, 2 => 1 per private subnet
resource "aws_nat_gateway" "nat" {
    count = 2

    allocation_id = aws_eip.nat_eip[count.index].id
    subnet_id = aws_subnet.public[count.index].id

    tags = {
        name = "${var.env}-aoc2024-natgw-${count.index}"
        environment = var.env
    }

    depends_on = [aws_internet_gateway.igw]
}

// route table for private subnet
resource "aws_route_table" "private" {
    count = 2
    vpc_id = aws_vpc.vpc.id 

    tags = {
        name = "${var.env}-aoc2024-private-rt-${count.index}"
        environment = var.env
    }
}

// 0.0.0.0 next-hop nat-gw
resource "aws_route" "private_gw" {
    count = 2

    route_table_id = aws_route_table.private[count.index].id
    destination_cidr_block = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat[count.index].id
}

// route table for public subnet
resource "aws_route_table" "public" {
    count = 2

    vpc_id = aws_vpc.vpc.id 
   
    tags = {
        name = "${var.env}-aoc2024-public-rt-${count.index}"
        environment = var.env
    }
}

// 0.0.0.0 next-hop internet-gw
resource "aws_route" "public_gw" {
    count = 2

    route_table_id = aws_route_table.public[count.index].id
    destination_cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
}

// associate route table with private subnets
resource "aws_route_table_association" "private" {
    count = 2

    subnet_id = aws_subnet.private[count.index].id
    route_table_id = aws_route_table.private[count.index].id
}

// associate route table with public subnets
resource "aws_route_table_association" "public" {
    count = 2

    subnet_id = aws_subnet.public[count.index].id
    route_table_id = aws_route_table.public[count.index].id
}