<template>
  <section id="cta" class="section-padding bg-gradient-to-br from-primary-600 to-cyan-600 relative overflow-hidden">
    <!-- Background Pattern -->
    <div class="absolute inset-0 opacity-10">
      <div class="absolute inset-0" :style="backgroundPattern"></div>
    </div>

    <!-- Floating Elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute top-20 left-10 w-32 h-32 bg-white/10 rounded-full animate-float"></div>
      <div class="absolute top-40 right-20 w-24 h-24 bg-white/10 rounded-full animate-float" style="animation-delay: 1s;"></div>
      <div class="absolute bottom-20 left-20 w-40 h-40 bg-white/10 rounded-full animate-float" style="animation-delay: 2s;"></div>
    </div>

    <div class="relative z-10 container-padding mx-auto text-center">
      <!-- Main CTA -->
      <h2 class="text-4xl md:text-5xl lg:text-6xl font-bold text-white mb-6 leading-tight">
        Become a Design Partner
      </h2>
      <p class="text-xl md:text-2xl text-white/90 mb-16 max-w-3xl mx-auto font-light leading-relaxed">
        Shape the future of explainable algorithms. Get early access, hands-on support, and special pricing as our founding partners.
      </p>

      <!-- Lead Capture Form -->
      <div class="max-w-lg mx-auto mb-16">
        <!-- Trust signals above form -->
        <div class="flex justify-center gap-4 mb-6">
          <div class="bg-white/10 backdrop-blur-sm rounded-lg px-4 py-2 flex items-center gap-2">
            <Shield class="w-4 h-4 text-white/80" />
            <span class="text-white/80 text-sm">DORA & US Compliant</span>
          </div>
          <div class="bg-white/10 backdrop-blur-sm rounded-lg px-4 py-2 flex items-center gap-2">
            <Lock class="w-4 h-4 text-white/80" />
            <span class="text-white/80 text-sm">SOC 2 Ready</span>
          </div>
          <div class="bg-white/10 backdrop-blur-sm rounded-lg px-4 py-2 flex items-center gap-2">
            <CheckCircle class="w-4 h-4 text-white/80" />
            <span class="text-white/80 text-sm">GDPR Compliant</span>
          </div>
        </div>

        <form @submit.prevent="handleSubmit" class="bg-white/10 backdrop-blur-xl rounded-3xl p-8 border border-white/30 premium-shadow-lg">
          <div class="space-y-4">
            <div>
              <input
                v-model="form.email"
                type="email"
                placeholder="Work Email *"
                class="w-full px-5 py-4 bg-white/20 border-2 border-white/30 rounded-xl text-white placeholder-white/60 focus:outline-none focus:border-white focus:bg-white/30 transition-all focus:ring-2 focus:ring-white/20 text-lg"
                required
              />
            </div>
            
            <div>
              <select
                v-model="form.role"
                class="w-full px-5 py-4 bg-white/20 border-2 border-white/30 rounded-xl text-white focus:outline-none focus:border-white focus:bg-white/30 transition-all focus:ring-2 focus:ring-white/20 text-lg"
                required
              >
                <option value="" class="text-gray-900">Your Role *</option>
                <option value="compliance" class="text-gray-900">Compliance Officer</option>
                <option value="security" class="text-gray-900">Security Professional</option>
                <option value="engineering" class="text-gray-900">Engineering Lead</option>
                <option value="executive" class="text-gray-900">C-Level Executive</option>
              </select>
            </div>
            
            <div>
              <select
                v-model="form.timeline"
                class="w-full px-5 py-4 bg-white/20 border-2 border-white/30 rounded-xl text-white focus:outline-none focus:border-white focus:bg-white/30 transition-all focus:ring-2 focus:ring-white/20 text-lg"
                required
              >
                <option value="" class="text-gray-900">Implementation Timeline *</option>
                <option value="immediate" class="text-gray-900">Immediate (Q1 2026)</option>
                <option value="quarter" class="text-gray-900">Next Quarter</option>
                <option value="half" class="text-gray-900">Next 6 Months</option>
              </select>
            </div>

            <button
              type="submit"
              :disabled="isSubmitting"
              class="w-full bg-gradient-to-r from-orange-500 to-red-500 text-white font-bold py-5 px-6 rounded-xl hover:from-orange-600 hover:to-red-600 transition-all transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 flex items-center justify-center gap-3 text-lg premium-shadow-lg mt-4"
            >
              <Sparkles v-if="!isSubmitting" class="w-6 h-6" />
              <Loader v-else class="w-6 h-6 animate-spin" />
              {{ isSubmitting ? 'Securing Your Spot...' : 'Get Compliant - Avoid Fines' }}
            </button>
          </div>

          <!-- Success Message -->
          <div v-if="showSuccess" class="mt-4 p-4 bg-green-500/20 border border-green-500/50 rounded-lg">
            <div class="flex items-center gap-2 text-green-100">
              <CheckCircle class="w-5 h-5" />
              <span class="font-medium">Thank you! We'll contact you within 24 hours to discuss the design partner program.</span>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="showError" class="mt-4 p-4 bg-red-500/20 border border-red-500/50 rounded-lg">
            <div class="flex items-center gap-2 text-red-100">
              <AlertCircle class="w-5 h-5" />
              <span class="font-medium">Something went wrong. Please try again.</span>
            </div>
          </div>
        </form>
      </div>

      <!-- Alternative CTAs -->
      <div class="flex flex-col sm:flex-row gap-6 justify-center items-center mb-16">
        <a
          href="tel:+1-555-DORA-HELP"
          class="flex items-center gap-3 px-6 py-3 bg-white/10 backdrop-blur-sm border border-white/30 rounded-lg text-white hover:bg-white/20 transition-colors"
        >
          <Phone class="w-5 h-5" />
          <span>Schedule a Call</span>
        </a>

        <a
          href="mailto:hunter@shannonlabs.dev"
          class="flex items-center gap-3 px-6 py-3 bg-white/10 backdrop-blur-sm border border-white/30 rounded-lg text-white hover:bg-white/20 transition-colors"
        >
          <Mail class="w-5 h-5" />
          <span>Email hunter@shannonlabs.dev</span>
        </a>
      </div>

      <!-- Trust Indicators -->
      <div class="border-t border-white/20 pt-12">
        <div class="grid md:grid-cols-2 gap-8 text-white">
          <div>
            <Shield class="w-12 h-12 mx-auto mb-4 text-white/80" />
            <h3 class="text-xl font-bold mb-2">Global Compliance</h3>
            <p class="text-white/70">Built for EU DORA, US regulations, and beyond</p>
          </div>
          <div>
            <Clock class="w-12 h-12 mx-auto mb-4 text-white/80" />
            <h3 class="text-xl font-bold mb-2">Quick Setup</h3>
            <p class="text-white/70">Go live in weeks, not months</p>
          </div>
        </div>
      </div>

      <!-- Final Urgency Message -->
      <div class="mt-16 text-center">
        <div class="inline-flex items-center gap-2 px-6 py-3 bg-red-500/20 backdrop-blur-sm rounded-full mb-6">
          <AlertTriangle class="w-5 h-5 text-red-300" />
          <span class="text-red-200 font-semibold">DORA Enforcement Active Since January 2025 | US Regulations Ongoing</span>
        </div>
        <p class="text-lg text-white/90 max-w-2xl mx-auto">
          Banks are being audited now. Every day without explainable fraud detection increases your compliance risk and potential fines.
        </p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Mail, Phone, Shield, Clock, AlertTriangle, CheckCircle, AlertCircle, Loader, Sparkles, Lock } from 'lucide-vue-next'

const isSubmitting = ref(false)
const showSuccess = ref(false)
const showError = ref(false)

const form = reactive({
  email: '',
  role: '',
  timeline: ''
})

const backgroundPattern = ref({
  backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
})

const handleSubmit = async () => {
  isSubmitting.value = true
  showError.value = false
  showSuccess.value = false

  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2000))

    // In a real implementation, you would send this data to your backend
    const formData = {
      ...form,
      timestamp: new Date().toISOString(),
      source: 'driftlock-landing-page'
    }

    console.log('Form submission:', formData)

    // Reset form
    Object.keys(form).forEach(key => {
      form[key as keyof typeof form] = ''
    })

    showSuccess.value = true

    // Hide success message after 10 seconds
    setTimeout(() => {
      showSuccess.value = false
    }, 10000)

  } catch (error) {
    console.error('Form submission error:', error)
    showError.value = true
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}
</style>