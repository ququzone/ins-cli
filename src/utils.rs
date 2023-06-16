use ethers::{
    abi::token::{LenientTokenizer, Tokenizer},
    types::U256,
};
use eyre::Result;
use std::str::FromStr;

pub fn strip_0x_prefix(s: &str) -> Result<String, &'static str> {
    Ok(s.strip_prefix("0x").unwrap_or(s).to_string())
}

pub fn parse_u256(s: &str) -> Result<U256> {
    Ok(if s.starts_with("0x") {
        U256::from_str(s)?
    } else {
        U256::from_dec_str(s)?
    })
}

pub fn parse_ether_value(value: &str) -> Result<U256> {
    Ok(if value.starts_with("0x") {
        U256::from_str(value)?
    } else {
        U256::from(LenientTokenizer::tokenize_uint(value)?)
    })
}
