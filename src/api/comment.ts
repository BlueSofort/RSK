import { api, userApi } from './client'

export interface CommentItem {
  id: number
  post_id: number
  user_id: number
  parent_id: number
  content: string
  status: string
  user_name: string
  user_avatar: string
  created_at: string
  updated_at: string
}

export interface CreateCommentPayload {
  post_id: number
  parent_id?: number
  content: string
}

export const commentAPI = {
  list: (postId: number, params?: any) =>
    api.get('/public/comments', { params: { post_id: postId, ...params } }),
  create: (data: CreateCommentPayload) =>
    userApi.post('/comments', data),
  delete: (id: number) =>
    userApi.delete(`/comments/${id}`),
}
