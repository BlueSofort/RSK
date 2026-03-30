import { userApi } from './client'
import type { ApiResponse } from './client'
import type { CreatePaymentPayload, PaymentCreateResult } from './types'

export const userOrderAPI = {
    preview: (data: any) => userApi.post<ApiResponse>('/orders/preview', data),
    create: (data: any) => userApi.post<ApiResponse>('/orders', data),
    createAndPay: (data: any) => userApi.post<ApiResponse>('/orders/create-and-pay', data),
    list: (params?: any) => userApi.get<ApiResponse>('/orders', { params }),
    detail: (orderNo: string, options?: any) => userApi.get<ApiResponse>(`/orders/${encodeURIComponent(orderNo)}`, options as any),
    cancel: (orderNo: string) => userApi.post<ApiResponse>(`/orders/${encodeURIComponent(orderNo)}/cancel`),
    downloadFulfillment: (orderNo: string) => userApi.get(`/orders/${encodeURIComponent(orderNo)}/fulfillment/download`, { responseType: 'blob' }),
}

export const guestOrderAPI = {
    preview: (data: any) => userApi.post<ApiResponse>('/guest/orders/preview', data),
    create: (data: any) => userApi.post<ApiResponse>('/guest/orders', data),
    createAndPay: (data: any) => userApi.post<ApiResponse>('/guest/orders/create-and-pay', data),
    list: (params: any) => userApi.get<ApiResponse>('/guest/orders', { params }),
    detail: (orderNo: string, params: any, options?: any) => userApi.get<ApiResponse>(`/guest/orders/${encodeURIComponent(orderNo)}`, { params, ...(options || {}) } as any),
    downloadFulfillment: (orderNo: string, params: any) => userApi.get(`/guest/orders/${encodeURIComponent(orderNo)}/fulfillment/download`, { params, responseType: 'blob' }),
    createPayment: (data: any) => userApi.post<ApiResponse>('/guest/payments', data),
    capturePayment: (id: number, data: any) => userApi.post<ApiResponse>(`/guest/payments/${id}/capture`, data),
    latestPayment: (params: any) => userApi.get<ApiResponse>('/guest/payments/latest', { params, silentBusinessError: true } as any),
}

export const paymentAPI = {
    create: (data: CreatePaymentPayload) => userApi.post<ApiResponse<PaymentCreateResult>>('/payments', data),
    capture: (id: number) => userApi.post<ApiResponse>(`/payments/${id}/capture`),
    latest: (params: any) => userApi.get<ApiResponse<PaymentCreateResult>>('/payments/latest', { params, silentBusinessError: true } as any),
}
