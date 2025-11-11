module infra {
    source = "../../../../modules/aws/infra"
    providers = {
        aws = aws
    }
    region = var.region
    env = var.env
    docker_image = var.docker_image
    alb_dns_name = var.alb_dns_name
    dns_zone = var.dns_zone
    alb_https = true
    alb_http = true 
    ecs_app_env_map = var.ecs_app_env_map
    ecs_app_env_map_secret = var.ecs_app_env_map_secret
}
