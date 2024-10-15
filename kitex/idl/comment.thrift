namespace go comment

// 引入用户相关定义
include "user.thrift"

// 发布or删除评论
struct CommentActionRequest {
  1: i64 video_id,           // 评论的视频id
  2: i32 action_type,        // 1-发布评论, 2-删除评论
  3: string comment_text,     // 用户填写的评论内容, action_type=1时使用
  4: i64 comment_id          // 要删除的评论id, action_type=2时使用
  5: string user_name     // 用户名
}

struct CommentActionResponse {
  1: i32 status_code,        // 状态码, 0成功, 其他值失败
  2: string status_msg,      // 返回状态描述
  3: Comment comment         // 评论成功返回评论内容, 不需要重新拉取整个评论列表
}

struct Comment {
  1: i64 id,                 // 评论id
  2: user.User user,         // 评论用户信息
  3: string content,         // 评论内容
  4: string create_date,     // 评论发布日期，格式mm-dd
}

// 评论列表
struct CommentListRequest {
  1: i64 video_id            // 评论的视频id
}

struct CommentListResponse {
  1: i32 status_code,        // 状态码
  2: string status_msg,      // 返回状态描述
  3: list<Comment> comment_list // 评论列表
}

service CommentService {
  CommentActionResponse CommentAction(1: CommentActionRequest req),
  CommentListResponse CommentList(1: CommentListRequest req)
}
