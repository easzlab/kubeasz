import { ElMessageBox, ElMessage } from 'element-plus'

let pendingOTP = null

export function useOTPDialog() {
  async function promptOTP() {
    if (pendingOTP) return pendingOTP

    try {
      const { value } = await ElMessageBox.prompt(
        'Please enter your OTP code for this high-risk operation',
        'OTP Verification',
        {
          confirmButtonText: 'Verify',
          cancelButtonText: 'Cancel',
          inputPattern: /^\d{6}$/,
          inputErrorMessage: 'OTP code must be 6 digits'
        }
      )
      pendingOTP = value
      return value
    } catch {
      return null
    }
  }

  function clearOTP() {
    pendingOTP = null
  }

  return { promptOTP, clearOTP }
}

export async function requestWithOTP(apiCall) {
  const { promptOTP, clearOTP } = useOTPDialog()

  try {
    return await apiCall()
  } catch (err) {
    if (err.response?.status === 403 && err.response?.data?.error === 'otp_required') {
      const code = await promptOTP()
      if (!code) {
        ElMessage.warning('OTP verification cancelled')
        throw err
      }
      // Retry with OTP code
      try {
        const result = await apiCall(code)
        clearOTP()
        return result
      } catch (retryErr) {
        clearOTP()
        if (retryErr.response?.data?.error === 'invalid_otp') {
          ElMessage.error('Invalid OTP code')
        }
        throw retryErr
      }
    }
    throw err
  }
}
