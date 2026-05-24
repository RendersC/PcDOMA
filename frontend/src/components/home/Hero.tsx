import { motion } from 'framer-motion';
import { Link } from 'react-router-dom';
import { ArrowRight, MapPin } from 'lucide-react';
import { TextScramble } from '../ui/text-scramble';
import { AnimatedGradientText } from '../ui/animated-gradient-text';
import { cn } from '../../lib/utils';

export function Hero() {
  return (
    // -mt-16 pulls the hero up behind the fixed navbar (main has pt-16).
    <section className="relative -mt-16 h-[100svh] min-h-[560px] w-full overflow-hidden">
      {/* Background video — full bleed */}
      <video
        autoPlay
        loop
        muted
        playsInline
        className="absolute inset-0 h-full w-full object-cover"
        src="/hero.mp4"
      />

      {/* Gradient overlay for legibility */}
      <div className="pointer-events-none absolute inset-0 bg-gradient-to-b from-gray-950/50 via-gray-950/25 to-gray-950/85" />

      {/* Content — bottom-aligned editorial layout */}
      <div className="absolute inset-x-0 bottom-0 px-5 pb-8 sm:px-8 md:px-12 md:pb-14">
        <div className="grid grid-cols-12 items-end gap-6">
          <div className="col-span-12 lg:col-span-8">
            <AnimatedGradientText className="mb-5">
              <MapPin className="mr-2 h-3.5 w-3.5 text-[#ffaa40]" />
              <span
                className={cn(
                  'inline animate-gradient bg-gradient-to-r from-[#ffaa40] via-[#9c40ff] to-[#ffaa40] bg-[length:var(--bg-size)_100%] bg-clip-text text-transparent',
                )}
              >
                Астана, Казахстан
              </span>
            </AnimatedGradientText>

            <TextScramble
              as="h1"
              duration={1.1}
              speed={0.05}
              className="font-extrabold leading-[0.9] tracking-[-0.04em] text-white text-[19vw] sm:text-[16vw] md:text-[13vw] lg:text-[11vw]"
            >
              PCDoma
            </TextScramble>
          </div>

          <div className="col-span-12 flex flex-col gap-5 pb-2 lg:col-span-4 lg:pb-4">
            <motion.p
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ duration: 0.8, delay: 0.5, ease: [0.16, 1, 0.3, 1] }}
              className="max-w-md text-sm leading-snug text-white/80 sm:text-base"
            >
              Арендуй мощный игровой или рабочий ПК с доставкой по городу. Собери конфигурацию
              сам или возьми готовый сетап — остальное сделает наш работник.
            </motion.p>

            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ duration: 0.8, delay: 0.7, ease: [0.16, 1, 0.3, 1] }}
              className="flex flex-wrap items-center gap-3"
            >
              <Link
                to="/catalog"
                className="group inline-flex items-center gap-2 rounded-full bg-white py-1.5 pl-5 pr-1.5 text-sm font-semibold text-gray-900 transition-all hover:gap-3"
              >
                Смотреть каталог
                <span className="flex h-8 w-8 items-center justify-center rounded-full bg-primary-600 transition-transform group-hover:scale-110">
                  <ArrowRight className="h-4 w-4 text-white" />
                </span>
              </Link>
              <Link
                to="/setups"
                className="rounded-full border border-white/30 bg-white/5 px-5 py-2.5 text-sm font-semibold text-white backdrop-blur-sm transition-colors hover:bg-white/15"
              >
                Готовые сетапы
              </Link>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
}
