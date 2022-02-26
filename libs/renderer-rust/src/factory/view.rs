use super::model::Contributor;
use svg::node::element::Circle;
use svg::node::element::Definitions;
use svg::node::element::Image;
use svg::node::element::Pattern;
use svg::node::element::SVG;
use svg::node::Text;
use svg::Document;

pub fn create_contributor_view(contributor: &Contributor, size: u16) -> SVG {
    let fill_id = format!("fill{}", contributor.id);
    let title = Text::new(format!("<title>{}</title>", contributor.login));

    let image = Image::new()
        .set("width", size)
        .set("height", size)
        .set("href", contributor.avatar_url.to_string());

    let fill = Pattern::new()
        .set("id", fill_id.clone())
        .set("x", 0)
        .set("y", 0)
        .set("width", size)
        .set("height", size)
        .set("patternUnits", "userSpaceOnUse")
        .add(image);

    let defs = Definitions::new().add(fill);

    let circle_border_width = 1;
    let circle = Circle::new()
        .set("cx", size as f32 / 2.0)
        .set("cy", size as f32 / 2.0)
        .set("r", size as f32 / 2.0)
        .set("stroke", "#c0c0c0")
        .set("stroke-width", circle_border_width)
        .set("fill", format!("url('#{}')", fill_id));

    Document::new()
        .set("width", size)
        .set("height", size)
        .add(title)
        .add(circle)
        .add(defs)
}

#[cfg(test)]
mod tests {
    use super::*;

    fn create_mock_contributor() -> Contributor {
        Contributor {
            id: 1,
            login: "login1".to_string(),
            avatar_url: "https://via.placeholder.com".to_string(),
            html_url: "htmlUrl1".to_string(),
            contributions: 1,
        }
    }

    #[test]
    fn test_render_title() {
        let result = create_contributor_view(&create_mock_contributor(), 64);
        println!("{}", result);
        assert!(result.to_string().contains("<title>login1</title>"));
    }

    #[test]
    fn test_contributor_view_is_sized() {
        let view = create_contributor_view(&create_mock_contributor(), 64);
        let attrs = view.get_inner().get_attributes();
        println!("{:?}", attrs);
        match (attrs.get("width"), attrs.get("height")) {
            (Some(width), Some(height)) => {
                assert_eq!(width.to_string(), "64");
                assert_eq!(height.to_string(), "64");
            }
            _ => assert!(false, "width not found"),
        }
    }
    #[test]
    fn test_contributor_view_is_a_circle() {
        let view = create_contributor_view(&create_mock_contributor(), 64);
        match view.get_inner().get_children().as_slice() {
            [_, shape, _] => {
                assert_eq!(
                    shape.to_string(),
                    "<circle cx=\"32\" cy=\"32\" fill=\"url('#fill1')\" r=\"32\" stroke=\"#c0c0c0\" stroke-width=\"1\"/>"
                );
            }
            _ => assert!(false, "circle not found"),
        }
    }
}
