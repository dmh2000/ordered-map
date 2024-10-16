use super::ordered_map::OrderedMap;

#[test]
fn test_new_ordered_map_empty() {
    let om: OrderedMap<i32, String> = OrderedMap::new();
    assert!(om.is_empty());
    assert_eq!(om.size(), 0);
}

#[test]
fn test_new_ordered_map_add_elements() {
    let mut om: OrderedMap<i32, String> = OrderedMap::new();
    om.put(1, "one".to_string());
    om.put(2, "two".to_string());
    om.put(3, "three".to_string());

    assert!(!om.is_empty());
    assert_eq!(om.size(), 3);

    let value = om.get(2);
    assert!(value.is_some());
    assert_eq!(value.unwrap(), "two");
}

#[test]
fn test_new_ordered_map_different_types() {
    let mut om_string: OrderedMap<String, i32> = OrderedMap::new();
    om_string.put("one".to_string(), 1);
    om_string.put("two".to_string(), 2);

    assert_eq!(om_string.size(), 2);

    let mut om_float: OrderedMap<f64, bool> = OrderedMap::new();
    om_float.put(1.1, true);
    om_float.put(2.2, false);

    assert_eq!(om_float.size(), 2);
}

#[test]
fn test_put_and_get() {
    let mut om: OrderedMap<i32, String> = OrderedMap::new();

    om.put(1, "one".to_string());
    om.put(2, "two".to_string());
    om.put(3, "three".to_string());

    let value = om.get(2);
    assert!(value.is_some());
    assert_eq!(value.unwrap(), "two");

    let value = om.get(4);
    assert!(value.is_none());

    om.put(2, "TWO".to_string());
    let value = om.get(2);
    assert!(value.is_some());
    assert_eq!(value.unwrap(), "TWO");
}

#[test]
fn test_delete() {
    let mut om: OrderedMap<i32, String> = OrderedMap::new();

    om.put(1, "one".to_string());
    om.put(2, "two".to_string());
    om.put(3, "three".to_string());

    om.delete(2);
    let value = om.get(2);
    assert!(value.is_none());

    assert_eq!(om.size(), 2);

    om.delete(4);
    assert_eq!(om.size(), 2);

    om.delete(1);
    om.delete(3);
    assert!(om.is_empty());
}

#[test]
fn test_put_and_get_with_different_types() {
    let mut om_string: OrderedMap<String, i32> = OrderedMap::new();
    om_string.put("one".to_string(), 1);
    om_string.put("two".to_string(), 2);

    let value = om_string.get("one".to_string());
    assert!(value.is_some());
    assert_eq!(*value.unwrap(), 1);

    let mut om_float: OrderedMap<f64, bool> = OrderedMap::new();
    om_float.put(1.1, true);
    om_float.put(2.2, false);

    let value = om_float.get(2.2);
    assert!(value.is_some());
    assert_eq!(*value.unwrap(), false);
}

#[test]
fn test_random_add_get_delete() {
    let mut om: OrderedMap<String, i32> = OrderedMap::new();

    let elements = vec![
        ("one".to_string(), 1),
        ("two".to_string(), 2),
        ("three".to_string(), 3),
    ];

    for (k, v) in &elements {
        om.put(k.clone(), *v);
    }

    assert_eq!(om.size(), elements.len());

    for (k, v) in &elements {
        let value = om.get(k.clone());
        assert!(value.is_some());
        assert_eq!(*value.unwrap(), *v);
    }

    for (k, _) in &elements {
        om.delete(k.clone());
    }

    assert!(om.is_empty());
}

#[test]
fn test_iterate_over_keys() {
    let mut om: OrderedMap<String, i32> = OrderedMap::new();

    let elements = vec![
        ("one".to_string(), 1),
        ("two".to_string(), 2),
        ("three".to_string(), 3),
    ];

    for (k, v) in &elements {
        om.put(k.clone(), *v);
    }

    assert_eq!(om.size(), elements.len());

    let keys = om.keys();
    for k in keys {
        let value = om.get(k.clone());
        assert!(value.is_some());
        assert_eq!(*value.unwrap(), *elements.iter().find(|(key, _)| key == k).unwrap().1);
    }

    for (k, _) in &elements {
        om.delete(k.clone());
    }

    assert!(om.is_empty());
}
