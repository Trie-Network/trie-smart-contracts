pub mod helpers;
pub mod imports;

use helpers::{call_mint_nft_api, MintNft, TransferFt, call_transfer_ft_api};

use rubixwasm_std::errors::WasmError;
use serde::{Deserialize, Serialize};
use rubixwasm_std::contract_fn;

#[derive(Serialize, Deserialize)]
pub struct PublishAssetReq {
    asset_id: String,
    asset_owner_did: String,
    asset_publish_description: String,
    asset_value: f64,
    asset_metadata: String,
    asset_filename: String,

    depin_provider_did: String,
    depin_hosting_cost: u32, // Hosting fees 

    ft_denom: String,
    ft_denom_creator: String
}

#[contract_fn]
pub fn publish_asset(publish_asset_req: PublishAssetReq) -> Result<String, WasmError> {
    let asset_creation_req = MintNft {
        did: publish_asset_req.asset_owner_did.clone(),
        nftId: publish_asset_req.asset_id.clone(),
        nftData: publish_asset_req.asset_publish_description,
        nftValue: publish_asset_req.asset_value,
        nftMetadata: publish_asset_req.asset_metadata,
        nftFilename: publish_asset_req.asset_filename.clone(),
    };

    let mint_nft_response = match call_mint_nft_api(asset_creation_req) {
        Ok(res) => {},
        Err(e) => return Err(WasmError::from(format!("failed while calling call_mint_nft_api, err: {:?}", e))),
    };

    // Pay Depin Provider in TRIE, and mention the NFT ID in the `comment`
    // for them to fetch NFT`
    let depin_payment_req = TransferFt {
        comment: format!("nft:{}", publish_asset_req.asset_id.clone()),
        ft_count: publish_asset_req.depin_hosting_cost as i32,
        ft_name: publish_asset_req.ft_denom,
        creatorDID: publish_asset_req.ft_denom_creator,
        sender: publish_asset_req.asset_owner_did,
        receiver: publish_asset_req.depin_provider_did.clone()
    };

    match call_transfer_ft_api(depin_payment_req) {
        Ok(_) => return Ok("".to_string()),
        Err(_) => return Err(WasmError { msg: format!("failed to send TRIE to DePin provider {}, please use 'resend_hosting_fees' contract function to retry sending TRIE tokens", publish_asset_req.depin_provider_did) }),
    };
}
