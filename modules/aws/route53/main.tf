data "aws_route53_zone" "primary" {
    name = var.dns_zone
}

resource "aws_route53_record" "alias" {
    zone_id = data.aws_route53_zone.primary.zone_id
    name = var.domain
    type = "CNAME"
    ttl = 300

    records = ["${var.alias}"]
}