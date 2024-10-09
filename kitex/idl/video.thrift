namespace go video

// 引入用户相关定义
include "user.thrift"

// feed视频流
struct Video {
  1: i64 id,  // 视频唯一标识
  2: user.User author,  // 视频作者信息
  3: string play_url,  // 视频播放地址
  4: string cover_url,  // 视频封面地址
  5: i64 favorite_count,  // 视频的点赞总数
  6: i64 comment_count,  // 视频的评论总数
  7: bool is_favorite,  // true-已点赞，false-未点赞
  8: string title,  // 视频标题
  9: i64 share_count  // 转发数量-本次暂不涉及
}

struct FeedRequest {
  1: i64 latest_time  // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
}

struct FeedResponse {
  1: i32 status_code,  // 状态码，0-成功，其他值-失败
  2: string status_msg,  // 返回状态描述
  3: list<Video> video_list,  // 视频列表
  4: i64 next_time  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

// 视频投稿
struct PublishActionRequest {
  1: binary data,  // 在 Thrift 中使用 binary 表示 bytes
  2: string title,
  3: string user_name
}

struct PublishActionResponse {
  1: i32 status_code,
  2: string status_msg
}

// 发布列表
struct PublishListRequest {
  1: i64 user_id
}

struct PublishListResponse {
  1: i32 status_code,
  2: string status_msg,
  3: list<Video> video_list  // 视频列表
}

service VideoService {
  FeedResponse Feed(1: FeedRequest req),
  PublishActionResponse PublishAction(1: PublishActionRequest req),
  PublishListResponse PublishList(1: PublishListRequest req)
}
