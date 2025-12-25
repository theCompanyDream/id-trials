import React from "react";

export default function About() {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 px-6 py-16">
      <div className="mx-auto max-w-4xl">
        <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-6">
          About This Project
        </h1>

        <p className="text-lg text-slate-300 leading-relaxed mb-10">
          This website is a practical playground for comparing modern identifier
          generation strategies. Instead of theory alone, it focuses on real,
          measurable behavior: speed, size, ordering, and operational tradeoffs
          you actually feel when systems scale.
        </p>

        <div className="grid gap-6 md:grid-cols-2 mb-12">
          <div className="rounded-2xl bg-slate-900/60 p-6 shadow-lg">
            <h2 className="text-xl font-semibold mb-2">What’s Being Tested</h2>
            <ul className="list-disc list-inside text-slate-300 space-y-1">
              <li>UUIDv4</li>
              <li>ULID</li>
              <li>KSUID</li>
              <li>Snowflake-style IDs</li>
              <li>NanoID</li>
              <li>CUID</li>
            </ul>
          </div>

          <div className="rounded-2xl bg-slate-900/60 p-6 shadow-lg">
            <h2 className="text-xl font-semibold mb-2">How They’re Evaluated</h2>
            <ul className="list-disc list-inside text-slate-300 space-y-1">
              <li>Generation performance</li>
              <li>Collision resistance</li>
              <li>Sortability and time ordering</li>
              <li>Storage footprint</li>
              <li>Distributed system behavior</li>
            </ul>
          </div>
        </div>

        <section className="mb-12">
          <h2 className="text-2xl font-semibold mb-4">Why Identifiers Matter</h2>
          <p className="text-slate-300 leading-relaxed">
            IDs are quiet infrastructure. They touch databases, indexes, caches,
            URLs, logs, and message queues. A good choice can reduce index
            fragmentation and improve query locality. A poor choice can quietly
            tax performance at scale. This project exists to make those tradeoffs
            visible, measurable, and easy to reason about.
          </p>
        </section>

        <section className="mb-12">
          <h2 className="text-2xl font-semibold mb-4">Goals of the Site</h2>
          <ul className="list-disc list-inside text-slate-300 space-y-2">
            <li>Provide side-by-side performance comparisons</li>
            <li>Highlight strengths and weaknesses of each ID type</li>
            <li>Show how ordering and randomness affect databases</li>
            <li>Offer practical guidance for real-world systems</li>
          </ul>
        </section>

        <div className="rounded-2xl bg-gradient-to-br from-indigo-500/20 to-cyan-500/20 p-6 shadow-lg">
          <h2 className="text-xl font-semibold mb-3">Who This Is For</h2>
          <p className="text-slate-200 leading-relaxed">
            This site is built for backend engineers, system designers, and
            anyone curious about how a single design decision can ripple through
            performance, scalability, and reliability. If you’ve ever wondered
            whether your IDs are helping or hurting, you’re in the right place.
          </p>
        </div>
      </div>
    </div>
  );
}
