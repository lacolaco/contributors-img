mod factory;

use self::factory::*;
use wasm_bindgen::prelude::*;

#[wasm_bindgen]
#[allow(clippy::boxed_local)]
pub fn render(contributors: Box<[JsValue]>, item_size: u16, max_columns: u16, gap: u16) -> String {
    factory::render(
        contributors
            .iter()
            .map(|item| item.into_serde().unwrap())
            .collect::<Vec<Contributor>>()
            .as_slice(),
        item_size,
        max_columns,
        gap,
    )
}
