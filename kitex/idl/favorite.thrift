namespace go favorite

// 引入视频相关定义
include "video.thrift"

// 点赞or取消点赞
struct FavoriteActionRequest {
  1: string user_name     // 用户名
  2: i64 video_id,     // 视频id
  3: i32 action_type   // 1-点赞，2-取消点赞
}

struct FavoriteActionResponse {
  1: i32 status_code,  // 状态码，0-成功，其他值-失败
  2: string status_msg  // 返回状态描述
}

// 点赞列表
struct FavoriteListRequest {
  1: i64 user_id      // 用户id
}

struct FavoriteListResponse {
  1: i32 status_code,  // 状态码，0-成功，其他值-失败
  2: string status_msg, // 返回状态描述
  3: list<video.Video> video_list  // 用户点赞视频列表
}

service FavoriteService {
  FavoriteActionResponse FavoriteAction(1: FavoriteActionRequest req),
  FavoriteListResponse FavoriteList(1: FavoriteListRequest req)
}
