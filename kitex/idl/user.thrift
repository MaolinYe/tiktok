namespace go user
// 用户注册
struct UserRegisterRequest {
  1: string username,  // 注册用户名，最长32个字符
  2: string password   // 密码，最长32个字符
}

struct UserRegisterResponse {
  1: i32 status_code,
  2: string status_msg,
  3: i64 user_id,
  4: string token
}

// 用户登录
struct UserLoginRequest {
  1: string username,
  2: string password
}

struct UserLoginResponse {
  1: i32 status_code,
  2: string status_msg,
  3: i64 user_id,
  4: string token
}

// 用户信息
struct User {
  1: i64 id,  // 用户id
  2: string name,  // 用户名称
  3: i64 follow_count,  // 关注总数
  4: i64 follower_count,  // 粉丝总数
  5: bool is_follow,  // true-已关注，false-未关注
  6: string avatar,  // 用户头像
  7: string background_image,  // 用户个人页顶部大图
  8: string signature,  // 个人简介
  9: i64 total_favorited,  // 获赞数量
  10: i64 work_count,  // 作品数量
  11: i64 favorite_count  // 点赞数量
}

struct UserInfoRequest {
  1: i64 user_id,
  2: string token
}

struct UserInfoResponse {
  1: i32 status_code,
  2: string status_msg,
  3: User user
}

service UserService {
  UserRegisterResponse Register(1: UserRegisterRequest req),
  UserLoginResponse Login(1: UserLoginRequest req),
  UserInfoResponse UserInfo(1: UserInfoRequest req)
}
