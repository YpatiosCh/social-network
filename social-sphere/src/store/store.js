import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export const useStore = create(
  persist(
    (set) => ({
      // State
      user: null,
      loading: false,

      // Manually set user data
      setUser: (userData) => {
        set({ user: userData })
      },

      // Clear user (on logout)
      clearUser: () => {
        set({ user: null })
      },
    }),
    {
      name: 'user',
      partialize: (state) => ({
        user: state.user
      }),
    }
  )
)