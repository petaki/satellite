import chai from 'chai';
import chaiAsPromised from 'chai-as-promised';
import electron from 'electron';
import path from 'path';
import { Application } from 'spectron';

chai.should();
chai.use(chaiAsPromised);

describe('Carrier', function() {
    this.timeout(10000);

    beforeEach(function() {
        this.app = new Application({
            args: [path.join(__dirname, '../build')],
            path: `${electron}`,
        });

        return this.app.start();
    });

    beforeEach(function() {
        chaiAsPromised.transferPromiseness = this.app.transferPromiseness;
    });

    afterEach(function() {
        if (this.app && this.app.isRunning()) {
            return this.app.stop();
        }
    });

    it('opens the window', function() {
        return this.app.client.waitUntilWindowLoaded()
            .getWindowCount().should.eventually.have.at.least(1)
            .browserWindow.isMinimized().should.eventually.be.false
            .browserWindow.isVisible().should.eventually.be.true
            .browserWindow.isFocused().should.eventually.be.true
            .browserWindow.getBounds().should.eventually.have.property('width').and.be.above(0)
            .browserWindow.getBounds().should.eventually.have.property('height').and.be.above(0);
    });
});
