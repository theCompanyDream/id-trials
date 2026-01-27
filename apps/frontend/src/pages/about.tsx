import { useUserStore } from '@backpack';

export default function About() {
  const idTypes = useUserStore((state) => state.idTypes)
  return (
    <main className="min-h-screen px-6 py-16 mx-auto max-w-4xl">
      <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-6">
        About This Project
      </h1>

      <p className="leading-relaxed mb-3">
        This website is a practical playground for comparing modern identifier
        generation strategies. Initially this was a user project that was a test
         and it was just a less functional explore page. Then I saw theo's <a href="https://www.youtube.com/watch?v=a-K2C3sf1_Q">video</a>
         on the problem with UUIDS. Then I had the idea why not build an entire site dedicated to exploring the performance of different Id types.
      </p>

      <section className="mb-3">
        <p className="leading-relaxed">
          IDs are quiet infrastructure. They touch everything from databases, indexes, caches, to
          URLs, logs, and message queues. I wanted to see if choosing a different id could yield to better results than other ids.
          This site benchmarks and compares popular ID generation strategies like UUID4, CUID, KSUID, ULID, Snowflake IDs, and Nano ID across various dimensions including performance, collision resistance, sortability, storage efficiency, and behavior in distributed systems.
        </p>
      </section>

      <section className="mb-3">
        <p className="leading-relaxed">
          This site is built for backend engineers, system designers, and
          anyone curious about how a single design decision can ripple through
          performance, scalability, and reliability. If you’ve ever wondered
          whether your IDs are helping or hurting, you’re in the right place.
        </p>
      </section>

      <section className="mb-3">
        <p className="leading-relaxed">
          It also this project is free and open source feel free to fork, star, and clone and run the test locally, or if you have ideas to improve the site or add more tests please open a PR on <a href="https://github.com/theCompanyDream/id-trials">GitHub</a>. This project was built in Golang and React and Typescript and deployed on Vercel.
        </p>
      </section>

      <section className="grid gap-6 md:grid-cols-2 mb-12 text-white">
        <div className="rounded-2xl bg-slate-900/60 p-6 shadow-lg">
          <h2 className="text-xl font-semibold mb-2">What’s Being Tested</h2>
          <ul className="list-disc list-inside space-y-1">
            {idTypes.map((idType) => (
              <li key={idType.value}>{idType.name}</li>
            ))}
          </ul>
        </div>

        <div className="rounded-2xl bg-slate-900/60 p-6 shadow-lg text-white">
          <h2 className="text-xl font-semibold mb-2">How They’re Evaluated</h2>
          <ul className="list-disc list-inside space-y-1">
            <li>Generation performance</li>
            <li>Collision resistance</li>
            <li>Sortability and time ordering</li>
            <li>Storage footprint</li>
            <li>Distributed system behavior</li>
          </ul>
        </div>
      </section>

      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">Goals of the Site</h2>
        <ul className="list-disc list-inside space-y-2">
          <li>Provide side by side performance comparisons</li>
          <li>Highlight strengths and weaknesses of each ID type</li>
          <li>Show how ordering and randomness affect databases</li>
          <li>Offer practical guidance for real-world systems</li>
        </ul>
      </section>
    </main>
  );
}
