import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Anomaly {
    id: string;
    stream_id: string;
    ncd: number;
    confidence: number;
    explanation: string;
    detected_at: string;
}

export interface ApiKey {
    id: string
    prefix: string
    fullKey: string
    status: 'Active' | 'Revoked'
    created: string
    createdDate: string
}

export const useAnomalyStore = defineStore('anomalies', () => {
    // State
    const anomalies = ref<Anomaly[]>([])
    const isConnected = ref(false)
    let eventSource: EventSource | null = null

    // Actions
    const connect = (apiKey: string, url: string = '/v1/stream/anomalies') => {
        if (eventSource) return // Already connected

        // In dev, we might need full URL if proxy isn't perfect, but relative should work with proxy
        eventSource = new EventSource(`${url}?api_key=${apiKey}`)

        eventSource.onopen = () => {
            isConnected.value = true
        }

        eventSource.onerror = () => {
            isConnected.value = false
        }

        eventSource.addEventListener('anomaly', (event) => {
            try {
                const anomaly = JSON.parse(event.data)
                addAnomaly(anomaly)
            } catch (e) {
                console.error('Failed to parse anomaly', e)
            }
        })
    }

    const disconnect = () => {
        if (eventSource) {
            eventSource.close()
            eventSource = null
            isConnected.value = false
        }
    }

    const addAnomaly = (anomaly: Anomaly) => {
        anomalies.value.unshift(anomaly)
        if (anomalies.value.length > 100) {
            anomalies.value.pop()
        }
    }

    return {
        anomalies,
        isConnected,
        connect,
        disconnect,
        addAnomaly
    }
})

export const useAuthStore = defineStore('auth', () => {
    const keys = ref<ApiKey[]>([
        {
            id: '1',
            prefix: 'dl_live_8a9b...',
            fullKey: 'dl_live_8a9b7c6d5e4f3g2h1i0j',
            status: 'Active',
            created: '2025-03-15T10:00:00Z',
            createdDate: 'March 15, 2025'
        }
    ])

    const generateKey = () => {
        const newId = Math.random().toString(36).substring(7)
        keys.value.push({
            id: newId,
            prefix: `dl_live_${newId}...`,
            fullKey: `dl_live_${newId}${Math.random().toString(36).substring(2)}`,
            status: 'Active',
            created: new Date().toISOString(),
            createdDate: new Date().toLocaleDateString()
        })
    }

    return {
        keys,
        generateKey
    }
})

