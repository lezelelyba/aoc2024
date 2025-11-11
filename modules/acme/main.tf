// register
resource "acme_registration" "reg" {
  email_address   = var.email
}

// generate key
resource "tls_private_key" "cert_private_key" {
  algorithm = var.algo
}

// generate request
resource "tls_cert_request" "req" {
  private_key_pem = tls_private_key.cert_private_key.private_key_pem

  subject {
    common_name = var.subject
  }
}

// request cert
resource "acme_certificate" "cert" {
    account_key_pem = acme_registration.reg.account_key_pem
    certificate_request_pem = tls_cert_request.req.cert_request_pem

    dns_challenge {
        provider = var.dns_provider
    }
}
