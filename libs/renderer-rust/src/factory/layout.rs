use super::model::Contributor;
use super::view::create_contributor_view;
use std::cmp::min;
use svg::node::element::SVG;
use svg::node::Node;

pub fn render(contributors: &[Contributor], item_size: u16, max_columns: u16, gap: u16) -> String {
    let columns = min(max_columns, contributors.len() as u16);
    let rows = (contributors.len() as f32 / columns as f32).ceil() as u16;

    let mut container = create_container(item_size, columns, rows, gap);
    for (i, c) in contributors.iter().enumerate() {
        let x = (i as f32 % columns as f32) as u16 * (item_size + gap);
        let y = (i as f32 / columns as f32).floor() as u16 * (item_size + gap);
        let view = create_contributor_view(c, item_size)
            .set("x", x)
            .set("y", y);
        container.append(view);
    }

    container.to_string()
}

fn create_container(item_size: u16, columns: u16, rows: u16, gap: u16) -> SVG {
    let width = item_size * columns + gap * (columns - 1);
    let height = item_size * rows + gap * (rows - 1);

    SVG::new()
        .set("width", width)
        .set("height", height)
        .set("viewBox", format!("0 0 {} {}", width, height))
}

#[test]
fn test_render_title_for_each_contributor() {
    let contributors = [
        Contributor {
            id: 1,
            login: "login1".to_string(),
            avatar_url: "https://via.placeholder.com".to_string(),
            html_url: "htmlUrl1".to_string(),
            contributions: 1,
        },
        Contributor {
            id: 2,
            login: "login2".to_string(),
            avatar_url: "https://via.placeholder.com".to_string(),
            html_url: "htmlUrl2".to_string(),
            contributions: 1,
        },
    ];

    let result = render(&contributors, 64, 12, 4);
    println!("{}", result);
    for login in ["login1", "login2"] {
        assert!(result.contains(login));
    }
}
