drop table if exists app_user, video, cart_item, search, search_result, job cascade;
drop index if exists app_user_token_index, cart_item_app_user_id, search_query;
