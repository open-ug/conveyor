use anyhow::{Result, anyhow};
use p256::ecdsa::SigningKey;
use pkcs8::DecodePrivateKey as _;
use rsa::{RsaPrivateKey, pkcs1::DecodeRsaPrivateKey, pkcs8::DecodePrivateKey};
use rustls_pemfile::read_all;
use std::io::Cursor;
use x509_parser::prelude::*;

#[derive(Debug)]
pub enum PrivateKey {
    Rsa(RsaPrivateKey),
    EcP256(SigningKey),
}

pub fn parse_cert_and_key(
    cert_pem: &[u8],
    key_pem: &[u8],
) -> Result<(Vec<X509Certificate<'static>>, PrivateKey)> {
    // Parse Private Key
    let mut key_reader = Cursor::new(key_pem);
    let items = read_all(&mut key_reader)?;

    if items.is_empty() {
        return Err(anyhow!("no private key PEM block found"));
    }

    let mut private_key: Option<PrivateKey> = None;

    for item in items {
        match item {
            rustls_pemfile::Item::Pkcs1Key(key) => {
                let rsa = RsaPrivateKey::from_pkcs1_der(&key)?;
                private_key = Some(PrivateKey::Rsa(rsa));
                break;
            }
            rustls_pemfile::Item::Pkcs8Key(key) => {
                // Try RSA first
                if let Ok(rsa) = RsaPrivateKey::from_pkcs8_der(&key) {
                    private_key = Some(PrivateKey::Rsa(rsa));
                    break;
                }

                // Try P256 EC
                if let Ok(ec) = SigningKey::from_pkcs8_der(&key) {
                    private_key = Some(PrivateKey::EcP256(ec));
                    break;
                }
            }
            rustls_pemfile::Item::Sec1Key(key) => {
                let ec = SigningKey::from_sec1_der(&key)?;
                private_key = Some(PrivateKey::EcP256(ec));
                break;
            }
            _ => {}
        }
    }

    let private_key = private_key.ok_or_else(|| anyhow!("unsupported private key type"))?;

    // Parse ALL Certificates
    let mut cert_reader = Cursor::new(cert_pem);
    let cert_items = read_all(&mut cert_reader)?;

    let mut certs = Vec::new();

    for item in cert_items {
        if let rustls_pemfile::Item::X509Certificate(cert_der) = item {
            let (_, cert) = X509Certificate::from_der(&cert_der)
                .map_err(|e| anyhow!("parse certificate in chain: {:?}", e))?;

            // Convert lifetime to 'static by owning the data
            let cert_owned = cert.to_owned();
            certs.push(cert_owned);
        }
    }

    if certs.is_empty() {
        return Err(anyhow!("no certificate PEM block found in cert file"));
    }

    // Validate Leaf Cert Matches Private Key
    let leaf = &certs[0];
    let cert_pub_key = leaf.public_key();

    let match_ok = match (&private_key, cert_pub_key.algorithm.algorithm.as_str()) {
        (PrivateKey::Rsa(rsa_key), _) => {
            if let Ok(cert_rsa) = rsa::RsaPublicKey::from_pkcs1_der(cert_pub_key.raw) {
                rsa_key.n() == cert_rsa.n() && rsa_key.e() == cert_rsa.e()
            } else {
                false
            }
        }
        (PrivateKey::EcP256(ec_key), _) => {
            if let Ok(cert_point) = p256::PublicKey::from_sec1_bytes(cert_pub_key.raw) {
                ec_key.verifying_key().to_encoded_point(false) == cert_point.to_encoded_point(false)
            } else {
                false
            }
        }
    };

    if !match_ok {
        return Err(anyhow!(
            "private key does not match first certificate in PEM chain"
        ));
    }

    Ok((certs, private_key))
}
